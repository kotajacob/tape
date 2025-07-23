package tape

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const EXT = ".mkv"

type tape struct {
	width, height int
	initialized   bool

	running  bool
	side     string
	frame    int
	duration time.Duration

	wg     *sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func New() tape {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	return tape{
		wg:     &wg,
		ctx:    ctx,
		cancel: cancel,
		side:   nextSide(),
	}
}

func nextSide() string {
	entries, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	names := make(map[string]struct{})
	for _, e := range entries {
		names[e.Name()] = struct{}{}
	}

	var name string
	for i := 0; true; i++ {
		r := getRune("", i)
		name = r + EXT
		if _, ok := names[name]; !ok {
			break
		}
	}
	return name
}

func getRune(prefix string, i int) string {
	if i/26 > 0 {
		prefix = getRune(prefix, (i/26)-1)
	}
	r := 'A' + rune(i%26)
	return prefix + string(r)
}

func (t tape) Init() tea.Cmd {
	return nil
}

func (t tape) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			if t.running {
				t.cancel()
				t.wg.Wait()
			}
			return t, tea.Quit
		case "enter", " ":
			if t.running {
				t.cancel()
				t.wg.Wait()
				t.ctx, t.cancel = context.WithCancel(context.Background())
				t.running = false
				t.frame = 0
				t.duration = 0
				t.side = nextSide()
			} else {
				t.wg.Add(1)
				go func() {
					wfRecorder(t.ctx, t.side)
					t.wg.Done()
				}()
				t.running = true
				cmds = append(cmds, t.tick())
			}
		}
	case tea.WindowSizeMsg:
		t.width = msg.Width
		t.height = msg.Height

		if !t.initialized {
			t.initialized = true
		}
	case tickMsg:
		if t.running {
			t.duration += time.Second
			t.frame = t.nextFrame()
			cmds = append(cmds, t.tick())
		}
	}

	return t, tea.Batch(cmds...)
}

// tickMsg is sent every second to update the run time display.
type tickMsg struct{}

// tick sends a TickMsg in one second.
func (t tape) tick() tea.Cmd {
	return tea.Tick(time.Second, func(_ time.Time) tea.Msg {
		return tickMsg{}
	})
}

func wfRecorder(ctx context.Context, filename string) {
	cmd := exec.CommandContext(
		ctx,
		"wf-recorder",
		"-c",
		"libx264",
		"-p",
		"crf=14",
		"-f",
		filename,
	)
	cmd.Cancel = func() error {
		return cmd.Process.Signal(syscall.SIGINT)
	}
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	cmd.Wait()
}

// nextFrame returns the int for the next frame.
func (t tape) nextFrame() int {
	f := t.frame + 1
	if f > 9 {
		f = 0
	}
	return f
}

func (t tape) View() string {
	if !t.initialized {
		return ""
	}
	return lipgloss.Place(
		t.width,
		t.height,
		lipgloss.Center,
		lipgloss.Center,
		t.tape(),
	)
}

// tape prints out the tape in its current state.
func (t tape) tape() string {
	var s strings.Builder
	dots := drawFrame(t.frame)
	side := drawSide(t.side)

	s.WriteString(`.------------------------.`)
	s.WriteString("\n")
	s.WriteString(`|\\////////      `)
	s.WriteString(mins(t.duration))
	s.WriteString(` min |`)
	s.WriteString("\n")
	s.WriteString(`| \/  __  ______  __     |`)
	s.WriteString("\n")
	s.WriteString(`|    /  \|\`)
	s.WriteString(dots[0])
	s.WriteString(`|/  \    |`)
	s.WriteString("\n")
	s.WriteString(`|    \  /|/`)
	s.WriteString(dots[1])
	s.WriteString(`|\  /    |`)
	s.WriteString("\n")
	s.WriteString(`| `)
	s.WriteString(side)
	s.WriteString(`‾‾  ‾‾‾‾‾‾  ‾‾     |`)
	s.WriteString("\n")
	s.WriteString(`|    ________________    |`)
	s.WriteString("\n")
	s.WriteString(`|___/_._o________o_._\___|`)
	return s.String()
}

// mins renders a 3 digit number of minutes from a duration.
func mins(d time.Duration) string {
	if d.Minutes() > 999 {
		d = time.Minute * 999
	}
	return fmt.Sprintf("%3.0f", d.Minutes())
}

// drawFrame renders the top followed by bottom frames for the tape.
func drawFrame(f int) [2]string {
	switch f {
	default:
		return [2]string{"·····", "     "}
	case 1:
		return [2]string{" ····", "    ·"}
	case 2:
		return [2]string{"  ···", "   ··"}
	case 3:
		return [2]string{"   ··", "  ···"}
	case 4:
		return [2]string{"    ·", " ····"}
	case 5:
		return [2]string{"     ", "·····"}
	case 6:
		return [2]string{"·    ", "···· "}
	case 7:
		return [2]string{"··   ", "···  "}
	case 8:
		return [2]string{"···  ", "··   "}
	case 9:
		return [2]string{"···· ", "·    "}
	}
}

func drawSide(s string) string {
	side := strings.TrimSuffix(s, EXT)
	if len(side) > 4 {
		return "ZZZZ"
	}
	for len(side) < 4 {
		side = side + " "
	}
	return side
}
