package main

import (
    "github.com/andrewmaach/GreedyAI/greedy"
    "math/rand"
    "time"
    "fmt"
    "strconv"
    "sync"
)

const matchSize int = 5

func main() {
    rand.Seed(time.Now().UnixNano())
    
    
    
    pool := SimpleStrategyGenerator("Adam", 9998)
   
    pool = append(pool, greedy.Strategy(SimpleStrategy{
                "R2-D2",
                [8]int{1500, 800, 700, 500, 500,500,400,0},
                [8]int{500, 700, 1000, 1500, 2000,2500,3000,10000},
                }))
            
    pool = append(pool, greedy.Strategy(SimpleStrategy{
                "BB-8",
                [8]int{800, 800, 700, 500, 500,500,400,0},
                [8]int{100, 100, 100, 100, 100,100,100,1000},
            }))
    
    
    records := War(pool)
    
    for _, item := range records {
        if item.wins > 0 {
            fmt.Printf("%s won %d of %d games.\n", item.id, item.wins, item.matches)
        }
    }
}


func War(pool []greedy.Strategy) map[string]WinLossRecord {
    
    battleResults := make(chan WinLossRecord, 100)
    battlesToFight := make(chan []greedy.Strategy, 100)
    
    go func(){
        for i := 0; i < 1000000; i++ {
           battlesToFight <- CreateBattle(pool, battleResults)
        }
        close(battlesToFight)
    }()
    
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(){
            defer wg.Done()
            for battle := range battlesToFight {
                winner := RepeatMatches(battle)
                battleResults <- WinLossRecord{
                    id: winner,
                    wins: 1,
                }
            }

        }()
    }
    
    go func() {
        wg.Wait()
        close(battleResults)
    }()
    
    records := make(map[string]WinLossRecord)
    
    for battle := range battleResults {
        if battle.wins > 0 {
            //fmt.Printf("%s won! \n", battle.id)
        }
        record, ok := records[battle.id]
        if !ok {
            records[battle.id] = battle
        } else {
            record.wins += battle.wins
            record.matches += battle.matches
            records[battle.id] = record
        }
    }
    
    return records
    
}

func CreateBattle(pool []greedy.Strategy, results chan WinLossRecord) []greedy.Strategy {
    strategies := []greedy.Strategy{}
    
    for c := 0; c < matchSize; c++ {
        item := pool[rand.Intn(len(pool))]
        
        results <- WinLossRecord{
            id: item.Id(),
            matches: 1,
        }
        
        strategies = append(strategies, item)
    }
    
    return strategies
}


// Returns the winner
func RepeatMatches(strategies []greedy.Strategy) string {
    game := greedy.CreateGame(strategies)
    game.Play()
    return game.Winner()
}

type WinLossRecord struct {
    id string
    wins int
    matches int
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
    return strategies
}

func Helpful() int {
    return rand.Intn(200) * 50
}