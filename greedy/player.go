package greedy

import (
    "math/rand"
    "fmt"
)

const FirstRoundMinScore int = 800
const WinningScore int = 10000
const debug = false

type Strategy interface {
    Id() string
    ShouldKeep(dice Dice, game *Game, minScore int) bool
    ShouldRoll(dice Dice, game *Game) bool
}

type Dice struct {
    Numbers [8]int
    Count int
    RunningScore int
}

var FreshDice Dice = Dice {
    Numbers: [8]int{0,0,0,0,0,0,0,0},
    Count: 8,
    RunningScore: 0,
}

type Player struct {
    Score int
    Plan Strategy
}

func CreateGame(strategies []Strategy) Game {
    game := Game{[]*Player{}, false}
    for _, strategy := range strategies {
        game.Players = append(game.Players, &Player{
            Score: 0,
            Plan: strategy,
        })
    }
    return game
}

type Game struct {
    Players []*Player
    LastRound bool
}

func (g *Game) Winner() string {
    highest := -1
    id := ""
    for _, player := range g.Players {
        if player.Score > highest {
            highest = player.Score
            id = player.Plan.Id()
        }
    }
    
    return id
}

func (g *Game) HighestScore() int {
    highest := 0
    for _, player := range g.Players {
        if player.Score > highest {
            highest = player.Score
        }
    }
    
    return highest
}

func (g *Game) Play() {
    turns := make(chan *Player, len(g.Players) + 1)
    
    for _, player := range g.Players {
        turns <- player
    }
    
    dice := FreshDice
    
    turnsLeft := 100000
    
    for player := range turns {
        minScore := 0
        
        // First round minimum score.
        if player.Score == 0 {
            minScore = FirstRoundMinScore
        }
        
        if g.LastRound {
            minScore = g.HighestScore() - player.Score + 1
        }
        
        dice = player.RunTurn(dice, minScore, g)
        if debug{fmt.Printf("%s scored %d, totaling %d \n \n", player.Plan.Id(), dice.RunningScore, player.Score)}
        

        
        if !g.LastRound {
            if player.Score >= WinningScore {
                g.LastRound = true
                close(turns)
            } else {
                 turns <- player
            }
           
        }
        
        turnsLeft -= 1
        if turnsLeft < 0 {
            panic("Too many turns in game.")
        }
    }
}

func (p *Player) RunTurn(passedDice Dice, minScore int, game *Game) Dice {
    dice := FreshDice
    kept := false
    if p.Plan.ShouldKeep(passedDice, game, minScore) {
        dice = passedDice
        if debug{fmt.Println("Keeping dice")}
        kept = true
    }
    

    
    for kept || p.Plan.ShouldRoll(dice, game) || dice.RunningScore < minScore {
        if !dice.Roll() {
            return FreshDice
        }
        kept = false
    }
    
    p.Score += dice.RunningScore
    
    return dice
}

func (d *Dice) Roll() bool {
    for i := 0; i < d.Count; i++ {
        d.Numbers[i] = rand.Intn(6) + 1
    }
    
    return d.evaluateScore()
}

func (d *Dice) evaluateScore() bool {
    rollScore := 0
    
    type Evaluation struct {
        Count [6]int
    }
    
     if debug{fmt.Printf("rolled: ")}
    
     // Count how many there are of each number.
    var eval Evaluation
    for i := 0; i < d.Count; i++ {
        if debug{fmt.Printf("%d ", d.Numbers[i])}
        eval.Count[d.Numbers[i] - 1] += 1
    }
    
    
    
    // Check for numbers with 3 or more.
    for i := 0; i < 6; i++ {
         if eval.Count[i] >= 3 {
             eval.Count[i] -= 3 // Decrement evaluation numbers.
             d.Count -= 3 // Decrement dice count.
             if i == 0 {
                 rollScore += 1000
             } else {
                 rollScore += (i + 1) * 100
             }
             i -= 1 // To catch two sets of the same number.
         }
         
    }
    
    // Check for ones.
     rollScore += eval.Count[0] * 100
     d.Count -= eval.Count[0] // Decrement dice count.
     
    // Check for fives.
     rollScore += eval.Count[4] * 50
     d.Count -= eval.Count[4] // Decrement dice count.
     
     
     if d.Count == 0 {
         d.Count = 8
     }
     
     if debug{fmt.Printf("(%d points, %d die left, total %d)\n", rollScore, d.Count, d.RunningScore)}
     
     if rollScore == 0 {
         d.RunningScore = 0
     } else {
         d.RunningScore += rollScore
     }
     return rollScore != 0
     
}

func (d *Dice) Print() {
    if debug{fmt.Printf("   dice left: %d, running score: %d, dice:", d.Count, d.RunningScore)}
    for i := 0; i < 8; i++ {
        if debug{fmt.Printf("%d ", d.Numbers[i])}
    }
    if debug{fmt.Println()}
}