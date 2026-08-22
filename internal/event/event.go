package event
import ("bytes"; "fmt"; "strconv")
type Event struct {
        ID int64
        Name string
        Data []byte
}
func Encode(ev Event) []byte {
        var b bytes.Buffer
        if ev.Name != "" {
                fmt.Fprintf(&b, "event: %s\n", ev.Name)
        }
        if ev.ID > 0 {
                fmt.Fprintf(&b, "id: %s\n", strconv.FormatInt(ev.ID, 10))
        }
        for _, line := range bytes.Split(ev.Data, []byte{'\n'}) {
                b.WriteString("data: ")
                b.Write(line)
                b.WriteByte('\n')
        }
        b.WriteByte('\n')
        return b.Bytes()
}
type Encoder interface { Encode(Event) []byte }
type DefaultEncoder struct{}
func (DefaultEncoder) Encode(ev Event) []byte { return Encode(ev) }
