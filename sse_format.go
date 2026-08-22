package ssehub
import ("bytes"; "fmt"; "strconv")
func FormatSSE(event, data string, id int) []byte {
        var b bytes.Buffer
        if event != "" { fmt.Fprintf(&b, "event: %s\n", event) }
        if id > 0 { fmt.Fprintf(&b, "id: %s\n", strconv.Itoa(id)) }
        for _, line := range bytes.Split([]byte(data), []byte("\n")) {
                b.WriteString("data: ")
                b.Write(line)
                b.WriteByte('\n')
        }
        b.WriteByte('\n')
        return b.Bytes()
}
