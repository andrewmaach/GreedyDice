package main

import (
    "github.com/andrewmaach/GreedyAI/greedy"
    "math/rand"
    "time"
    "fmt"
)

func main() {
    rand.Seed(time.Now().UnixNano())
    
    for i := 0 ; i < 1 ; i++ {
        strategySet := []greedy.Strategy{
            greedy.Strategy{"R2-D2", 200},
            greedy.Strategy{"BB-8", 400},
        }
        
        game := greedy.CreateGame(strategySet)
        
        game.Play()
        for _, player := range game.Players {
            fmt.Printf("%s: %d\n", player.Plan.Id, player.Score)
        }
    }
    
}