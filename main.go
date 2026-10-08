package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 480
	screenHeight = 640

	basketWidth  = 80
	basketHeight = 20
	basketMargin = 10 // 篮筐距底部距离

	ballRadius = 10
)

// Ball 下落的小球
type Ball struct {
	x, y  float64
	speed float64 // 每帧下落像素
	color color.Color
}

// Game 游戏状态
type Game struct {
	balls       []Ball
	basketX     float64 // 篮筐左上角 x
	score       int
	missed      int
	spawnTicker int // 距离下次发射还剩的帧数
	gameOver    bool
}

var ballColors = []color.Color{
	color.RGBA{255, 99, 71, 255},  // 番茄红
	color.RGBA{255, 215, 0, 255},  // 金黄
	color.RGBA{50, 205, 50, 255},  // 绿色
	color.RGBA{135, 206, 250, 255}, // 天蓝
	color.RGBA{221, 160, 221, 255}, // 紫色
}

// randomSpawnInterval 返回 1-3 秒对应的帧数（60帧/秒）
func randomSpawnInterval() int {
	return 60 + rand.IntN(121) // 60~180 帧 = 1~3 秒
}

func (g *Game) spawnBall() {
	g.balls = append(g.balls, Ball{
		x:     ballRadius + rand.Float64()*(screenWidth-2*ballRadius),
		y:     -ballRadius,
		speed: 2 + rand.Float64()*2, // 2~4 像素/帧，速度适中可接
		color: ballColors[rand.IntN(len(ballColors))],
	})
}

func (g *Game) Update() error {
	if g.gameOver {
		// 按空格重新开始
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.reset()
		}
		return nil
	}

	// 篮筐跟随鼠标水平移动
	mx, _ := ebiten.CursorPosition()
	g.basketX = float64(mx) - basketWidth/2
	if g.basketX < 0 {
		g.basketX = 0
	}
	if g.basketX > screenWidth-basketWidth {
		g.basketX = screenWidth - basketWidth
	}

	// 发射小球
	g.spawnTicker--
	if g.spawnTicker <= 0 {
		g.spawnBall()
		g.spawnTicker = randomSpawnInterval()
	}

	// 更新小球位置并检测碰撞
	basketY := float64(screenHeight - basketMargin - basketHeight)
	alive := g.balls[:0]
	for _, b := range g.balls {
		b.y += b.speed

		// 小球到达篮筐高度：判断是否接住
		if b.y+ballRadius >= basketY && b.y-ballRadius <= basketY+basketHeight {
			if b.x >= g.basketX && b.x <= g.basketX+basketWidth {
				g.score++
				continue // 接住，移除小球
			}
		}

		// 掉出屏幕：漏接
		if b.y-ballRadius > screenHeight {
			g.missed++
			if g.missed >= 10 {
				g.gameOver = true
			}
			continue
		}

		alive = append(alive, b)
	}
	g.balls = alive

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 30, 46, 255})

	// 画小球
	for _, b := range g.balls {
		vector.DrawFilledCircle(screen, float32(b.x), float32(b.y), ballRadius, b.color, true)
	}

	// 画篮筐（跟随鼠标）
	basketY := float32(screenHeight - basketMargin - basketHeight)
	vector.DrawFilledRect(screen, float32(g.basketX), basketY, basketWidth, basketHeight, color.RGBA{255, 165, 0, 255}, true)
	// 篮筐边缘加亮，看起来像个箩筐
	vector.StrokeRect(screen, float32(g.basketX), basketY, basketWidth, basketHeight, 3, color.RGBA{255, 200, 60, 255}, true)

	// 分数
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Score: %d   Missed: %d/10", g.score, g.missed), 10, 10)

	if g.gameOver {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("GAME OVER! Final Score: %d", g.score), screenWidth/2-110, screenHeight/2-10)
		ebitenutil.DebugPrintAt(screen, "Press SPACE to restart", screenWidth/2-80, screenHeight/2+15)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) reset() {
	g.balls = nil
	g.score = 0
	g.missed = 0
	g.spawnTicker = randomSpawnInterval()
	g.gameOver = false
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Catch the Balls - 接球小游戏")
	g := &Game{}
	g.reset()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
