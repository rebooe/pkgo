package middle

import (
	"net/http"
)

// 处理跨域中间件
func CorsHandler(AllowOrigin string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", AllowOrigin) // 可将将 * 替换为指定的域名
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
}
