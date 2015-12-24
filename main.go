package main

import (
    "github.com/andrewmaach/GreedyAI/greedy"
    "math/rand"
    "time"
    "fmt"
    "strconv"
)

func main() {
    rand.Seed(time.Now().UnixNano())
    
    for i := 0 ; i < 10 ; i++ {
        strategySet := []greedy.Strategy{
            greedy.Strategy(SimpleStrategy{
                "R2-D2",
                [8]int{1500, 800, 700, 500, 500,500,400,0},
                [8]int{500, 700, 1000, 1500, 2000,2500,3000,10000},
                }),
            greedy.Strategy(SimpleStrategy{
                "BB-8",
                [8]int{800, 800, 700, 500, 500,500,400,0},
                [8]int{100, 100, 100, 100, 100,100,100,1000},
            }),
        }
        
        game := greedy.CreateGame(strategySet)
        
        game.Play()
        for _, player := range game.Players {
            fmt.Printf("%s: %d\n", player.Plan.Id(), player.Score)
        }
    }
    
}

type SimpleStrategy struct {
    id string
    MinKeepPoints [8]int
    MinStopPoints [8]int
}

func (s SimpleStrategy) Id()  string {
    return s.id
}


func (s SimpleStrategy) ShouldKeep(dice greedy.Dice, game *greedy.Game, minScore int)  bool {
    return dice.RunningScore > s.MinKeepPoints[dice.Count - 1]
}

func (s SimpleStrategy) ShouldRoll(dice greedy.Dice, game *greedy.Game)  bool {
    return dice.RunningScore < s.MinStopPoints[dice.Count - 1]
}

func SimpleStrategyGenerator(prefix string, count int) []greedy.Strategy {
    strategies := []greedy.Strategy{}
    for i := 0; i < count; i++ {
        strategies = append(strategies, greedy.Strategy(SimpleStrategy{
                prefix+":"+strconv.Itoa(i+1),
                [8]int{Helpful(), Helpful(), Helpful(), Helpful(), Helpful(),Helpful(),Helpful(),Helpful()},
                [8]int{Helpful(), Helpful(), Helpful(), Helpful(), Helpful(),Helpful(),Helpful(),Helpful()},
            }))
    }
    return nil
}

func Helpful() int {
    return rand.Intn(200) * 50
}