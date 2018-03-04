package sidedness

import (
	"testing"
	"github.com/franela/goblin"
	"github.com/intdxdt/geom"
	"time"
)

func TestCmp(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("homotopic sidedness test", func() {
		g.It("should test sidedness to a context geometry", func() {
			g.Timeout(1 * time.Hour)
			var cwkt = "POLYGON (( 278 307, 270 298, 274 286, 279 272, 301 274, 308 288, 311 304, 296 308, 278 307 ))"
			var wkt = "LINESTRING ( 155 171, 207 166, 253 175, 317 171, 367 182, 400 200, 428 249, 417 291, 383 324, 361 332, 333 347, 314 357, 257 383, 204 370, 176 337, 180 305, 214 295, 244 302, 281 332, 316 328, 331 306, 332 291, 315 265, 285 250, 247 261, 231 276, 195 264, 187 230, 216 215, 257 226, 273 217, 273 205, 240 197, 200 200, 178 193, 157 226, 156 246, 151 263, 120 264, 95 249, 89 261, 100 300, 116 359, 139 389, 172 413, 211 425, 256 430, 289 431, 348 427 )"
			var coords = geom.NewLineStringFromWKT(wkt).Coordinates()
			g.Assert(
				IsHomotopic(coords, geom.NewPolygonFromWKT(cwkt)),
			).IsTrue()
			g.Assert(
				Homotopy(coords, geom.NewPolygonFromWKT(cwkt)),
			).IsTrue()

			cwkt = "POLYGON (( 221 347, 205 334, 221 322, 234 324, 237 342, 221 347 ))"
			g.Assert(
				IsHomotopic(coords, geom.NewPolygonFromWKT(cwkt)),
			).IsFalse()
			g.Assert(
				Homotopy(coords, geom.NewPolygonFromWKT(cwkt)),
			).IsFalse()

		})

	})
}
