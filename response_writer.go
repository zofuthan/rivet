package rivet

import (
	"net/http"
)

/**
ResponseWriter 提供状态支持
*/
type ResponseWriter interface {
	http.ResponseWriter
	http.Flusher
	// Status 返回调用 WriteHeader 设定的值, 初始值为 0
	Status() int
	// Size 返回调用 Write 写入的字节数, 初始值为 0
	Size() int

	/**
	Written 返回是否已经写入了内容.
	包括两种情况, WriteHeader 和 Write. 实现是如何判断的可能有差异.
	*/
	Written() bool
}

/**
NewResponseWriterFakeFlusher 返回 ResponseWriter 实例, 可能是伪 http.Flusher.
如果 rw 已经实现了 ResponseWriter 接口, 返回 rw.(ResponseWriter).
否则返回 &ResponseWriteFakeFlusher 伪 http.Flusher 实例.
*/
func NewResponseWriterFakeFlusher(rw http.ResponseWriter) ResponseWriter {

	if res, ok := rw.(ResponseWriter); ok {
		return res
	}
	return &ResponseWriteFakeFlusher{rw, 0, 0}
}

/**
ResponseWriteFakeFlusher 实现了 http.ResponseWriter 接口和伪 http.Flusher 接口.
Flush() 是个方法, 是否支持 Flusher 取决于原 http.ResponseWriter 实例.
*/
type ResponseWriteFakeFlusher struct {
	http.ResponseWriter
	status int
	size   int
}

// Flush() 是个伪方法, 是否支持 Flusher 取决于原 http.ResponseWriter 实例.
func (rw *ResponseWriteFakeFlusher) Flush() {
	flusher, ok := rw.ResponseWriter.(http.Flusher)
	if ok {
		flusher.Flush()
	}
}

/**
WriteHeader 向相应发送状态码 s.
*/
func (rw *ResponseWriteFakeFlusher) WriteHeader(s int) {
	rw.ResponseWriter.WriteHeader(s)
	rw.status = s
}

/**
Write 向相应写入 b, 返回本次写入的字节和发生的错误.
*/
func (rw *ResponseWriteFakeFlusher) Write(b []byte) (int, error) {
	if !rw.Written() {
		rw.WriteHeader(http.StatusOK)
	}
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

/**
Status 返回通过 WriteHeader 设置的值.
*/
func (rw *ResponseWriteFakeFlusher) Status() int {
	return rw.status
}

/**
Size 返回通过 Write 的总字节数.
*/
func (rw *ResponseWriteFakeFlusher) Size() int {
	return rw.size
}

/**
Written 返回  Status()!=0 || Size()!=0 的结果
*/
func (rw *ResponseWriteFakeFlusher) Written() bool {
	return rw.status != 0 || rw.size != 0
}
