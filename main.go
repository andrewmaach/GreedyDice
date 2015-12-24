package main

import (
    "github.com/andrewmaach/GreedyAI/greedy"
    "math/rand"
    "time"
    "fmt"
    "strconv"
)

const matchSize int = 5

func main() {
    rand.Seed(time.Now().UnixNano())
    
    pool := SimpleStrategyGenerator("Adam", 1000)
   
    pool = append(pool, Recordify(greedy.Strategy(SimpleStrategy{
                "R2-D2",
                [8]int{1500, 800, 700, 500, 500,500,400,0},
                [8]int{500, 700, 1000, 1500, 2000,2500,3000,10000},
                })))
            
    pool = append(pool, Recordify(greedy.Strategy(SimpleStrategy{
                "BB-8",
                [8]int{800, 800, 700, 500, 500,500,400,0},
                [8]int{100, 100, 100, 100, 100,100,100,1000},
            })))
    
    
    
    for i := 0 ; i < 10 ; i++ {
        
        chosen := []*WinLossRecord{}
        strategies := []greedy.Strategy{}
        
        for c := 0; c < matchSize; c++ {
            item := pool[rand.Intn(len(pool))]
            chosen = append(chosen, item)
            strategies = append(strategies, item.strategy)
        }
        
        
        game := greedy.CreateGame(strategies)
        game.Play()
        for _, player := range game.Players {
            fmt.Printf("%s: %d\n", player.Plan.Id(), player.Score)
        }
    }
    
}

type WinLossRecord struct {
    strategy greedy.Strategy
    winavg float64
    matches int
}

func Recordify(strategy greedy.Strategy) *WinLossRecord {

    v := new(WinLossRecord)
    v.strategy = strategy
    v.winavg = 0.0
    v.matches = 0
    return v
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

func SimpleStrategyGenerator(prefix string, count int) []*WinLossRecord {
    strategies := []*WinLossRecord{}
    for i := 0; i < count; i++ {
        strategies = append(strategies, Recordify(greedy.Strategy(SimpleStrategy{
                prefix+":"+strconv.Itoa(i+1),
                [8]int{Helpful(), Helpful(), Helpful(), Helpful(), Helpful(),Helpful(),Helpful(),Helpful()},
                [8]int{Helpful(), Helpful(), Helpful(), Helpful(), Helpful(),Helpful(),Helpful(),Helpful()},
            })))
    }
    return strategies
}

func Helpful() int {
    return rand.Intn(200) * 50
}