package engine

// type writer struct {
// 	conn   *Conn
// 	prefix string

// 	line struct {
// 		Kind MessageKind `json:"kind"`
// 		Data struct {
// 			Group string `json:"group"`
// 			Line  string `json:"line"`
// 		} `json:"data"`
// 	}
// }

// func NewLogWriter(c *Conn, prefix string) *LogWriter {
// 	w := new(LogWriter)
// 	w.conn = c
// 	w.prefix = prefix
// 	return w
// }

// func (w *LogWriter) Write(data []byte) (int, error) {
// 	// w.line.Data.Line = strings.TrimSuffix(string(data), "\n")
// 	w.line.Data.Line = string(data)
// 	w.line.Data.Group = w.prefix
// 	w.line.Kind = MessageLogs
// 	err := w.conn.Send(&w.line)
// 	return len(data), err
// }
