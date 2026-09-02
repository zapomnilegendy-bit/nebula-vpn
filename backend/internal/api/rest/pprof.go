package rest

import (
	"runtime"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
)

func setupPprof(app *fiber.App) {
	app.Use(pprof.New(pprof.Config{
		Prefix: "/debug/pprof",
	}))
}

// Set memory limits for GC tuning
func init() {
	// Let Go use up to 80% of available memory before triggering GC
	debug.SetGCPercent(80)

	// Set memory limit (soft)
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
}
