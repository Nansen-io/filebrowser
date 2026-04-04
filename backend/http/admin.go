package http

import (
	"net/http"
	"runtime"
	"time"

	"github.com/gtsteffaniak/filebrowser/backend/common/settings"
	"github.com/gtsteffaniak/filebrowser/backend/common/version"
	"github.com/gtsteffaniak/filebrowser/backend/indexing"
)

var serverStartTime = time.Now()

type adminStats struct {
	Uptime     int64                            `json:"uptime"`
	Version    string                           `json:"version"`
	Goroutines int                              `json:"goroutines"`
	MemAlloc   uint64                           `json:"memAlloc"`
	MemSys     uint64                           `json:"memSys"`
	NumGC      uint32                           `json:"numGC"`
	Sources    map[string]indexing.ReducedIndex `json:"sources"`
}

// adminStatsHandler returns runtime stats for the admin dashboard.
// @Summary Admin stats
// @Description Returns Go runtime stats and per-source index info for admins.
// @Tags Admin
// @Produce json
// @Success 200 {object} adminStats
// @Router /api/admin/stats [get]
func adminStatsHandler(w http.ResponseWriter, r *http.Request, d *requestContext) (int, error) {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	sources := settings.GetSources(d.user)
	sourceData := map[string]indexing.ReducedIndex{}
	for _, source := range sources {
		ri, err := indexing.GetIndexInfo(source, false)
		if err == nil {
			sourceData[source] = ri
		}
	}

	stats := adminStats{
		Uptime:     int64(time.Since(serverStartTime).Seconds()),
		Version:    version.Version,
		Goroutines: runtime.NumGoroutine(),
		MemAlloc:   ms.Alloc,
		MemSys:     ms.Sys,
		NumGC:      ms.NumGC,
		Sources:    sourceData,
	}

	return renderJSON(w, r, stats)
}
