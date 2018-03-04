package sidedness

import (
	"testing"
	"github.com/franela/goblin"
	"github.com/intdxdt/geom"
	"time"
)

func TestSidedness(t *testing.T) {
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

			cwkt = "POLYGON (( 221 347, 205 334, 221 322, 234 324, 237 342, 221 347 ))"
			wkt = "LINESTRING ( 155 171, 207 166, 253 175, 317 171, 366.3526316090967 196.8011373660995, 403.24889098928804 217.07380735521562, 449.06512516469047 258.02460073323016, 469.7432485535889 310.7335427049321, 431.6306289740506 348.8461622844704, 393.9234627942946 354.9279632812052, 371.21807240648457 309.92263590536743, 372.43443260583155 253.97006673540696, 415.4124929827577 249.91553273758373, 425.14337457753345 316.00443690210227, 404.465251188635 331.81711949361284, 359.0544704130149 294.10995331385686, 333 347, 365.94717820931436 412.50234605029505, 350.9454024173684 441.6949908346222, 269.44926906112164 441.6949908346222, 233.36391648049494 444.9386180328808, 205.38763189551466 402.77146445551926, 204 370, 176 337, 180 305, 214 295, 244 302, 281 332, 316 328, 331 306, 332 291, 315 265, 285 250, 247 261, 231 276, 195 264, 187 230, 216 215, 257 226, 273 217, 273 205, 240 197, 200 200, 178 193, 157 226, 156 246, 151 263, 120 264, 95 249, 89 261, 100 300, 116 359, 139 389, 172 413, 211 425, 256 430, 289 431, 348 427 )"
			coords = geom.NewLineStringFromWKT(wkt).Coordinates()
			g.Assert(
				Homotopy(coords, geom.NewPolygonFromWKT(cwkt)),
			).IsFalse()
			cwkt = "POLYGON (( 278 307, 270 298, 274 286, 279 272, 301 274, 308 288, 311 304, 296 308, 278 307 ))"
			g.Assert(
				Homotopy(coords, geom.NewPolygonFromWKT(cwkt)),
			).IsTrue()

			wkt = "LINESTRING ( 155 171, 207 166, 253 175, 317 171, 366.3526316090967 196.8011373660995, 403.24889098928804 217.07380735521562, 449.2399093677868 286.1290060157034, 451.783228702019 325.7276851722446, 415.8591984439269 355.71930167053, 418.6880722510337 402.11283210708103, 389.8335594185446 420.2176244725644, 350.9454024173684 441.6949908346222, 269.44926906112164 441.6949908346222, 233.36391648049494 444.9386180328808, 205.38763189551466 402.77146445551926, 204 370, 176 337, 180 305, 157 226, 95 249, 89 261, 100 300, 116 359, 139 389, 172 413, 211 425, 256 430, 289 431, 348 427 )"
			cwkt = "POLYGON (( 278 307, 270 298, 274 286, 279 272, 301 274, 308 288, 311 304, 296 308, 278 307 ))"
			coords = geom.NewLineStringFromWKT(wkt).Coordinates()
			g.Assert(
				Homotopy(coords, geom.NewPolygonFromWKT(cwkt)),
			).IsFalse()
			cwkt = "POLYGON (( 221 347, 205 334, 221 322, 234 324, 237 342, 221 347 ))"
			g.Assert(
				Homotopy(coords, geom.NewPolygonFromWKT(cwkt)),
			).IsTrue()

			coords = []*geom.Point {{2,2}, {5,2}, {7,2}, {9,2}}
			var lnconst = geom.NewLineString([]*geom.Point {{3,2}, {6,2}, {6.5,2}})
			g.Assert(
				Homotopy(coords, lnconst),
			).IsTrue()
			var plyconst = geom.NewPolygon([]*geom.Point {{4,1}, {4,3}, {5,3}, {5,1}})
			g.Assert(
				Homotopy(coords, plyconst),
			).IsTrue()
		})

	})
}
