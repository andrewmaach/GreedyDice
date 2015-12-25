package main

import (
    "github.com/andrewmaach/GreedyAI/greedy"
    "math/rand"
    "time"
    "fmt"
    "strconv"
    "sync"
    "sort"
)

const matchSize int = 5
const generations int = 100

const initialSeeds int = 20000
const keepHighestOf int = 10000
const breedCount int = 10000
const battlesPerGeneration int = 1000000
const generationCounter int = 1000

func main() {
    rand.Seed(time.Now().UnixNano())
    
    
    pool := GenerateSeeds(initialSeeds)
    for gen:=0; gen < generations; gen++ {
        fmt.Println()
        fmt.Println()
        fmt.Printf("War, generation #%d\n", gen)
        records := War(pool)
        
        // flatten records
        flatRecords := []WinLossRecord{}
        for _, v := range records {flatRecords = append(flatRecords, v)}
        
        // sort
        sort.Sort(ByWinPercent(flatRecords))
        
        topPrint := 1
        if gen == generations - 1 {topPrint = 10}
        for x:=0; x<topPrint; x++ {
            item := flatRecords[x]
             fmt.Printf("#%d: %s won %d of %d games.\n",x+1, item.id, item.wins, item.matches)
            GetStrategy(pool, item.id).Print()
        }
        
        // Keep the upper half
        
        highest := flatRecords[:keepHighestOf]
        
        nextGen := CollectFromRecords(highest, pool)
        
        fmt.Println("Breeding...")
        bred := 0
        for bred < breedCount {
            a := nextGen[rand.Intn(keepHighestOf)]
            b := nextGen[rand.Intn(keepHighestOf)]
            if a == b { continue }
            
            c := a.Breed(b, "bred:"+strconv.Itoa(bred+(gen*generationCounter)))
            if c == nil {continue}
            nextGen = append(nextGen, c)
            bred += 1
        }
        
        pool = nextGen
        
    }


}

func CollectFromRecords(records []WinLossRecord, oldPool []greedy.Strategy) []greedy.Strategy {
    pool := []greedy.Strategy{}
    for _, record := range records {
        e :=  GetStrategy(oldPool, record.id)
        if e != nil {
            pool = append(pool,e)
        }
        
    }
    return pool
}

func GetStrategy(pool []greedy.Strategy, id string) greedy.Strategy {
    for _, strategy := range pool {
        if strategy.Id() == id {return strategy}
    }
    return nil
}



type ByWinPercent []WinLossRecord

func (a ByWinPercent) Len() int           { return len(a) }
func (a ByWinPercent) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByWinPercent) Less(i, j int) bool {

    return a[i].Victory() > a[j].Victory()
}

func GenerateSeeds(count int)  []greedy.Strategy {
    pool := SimpleStrategyGenerator("Adam", count)
   
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
            
    return pool
}




func War(pool []greedy.Strategy) map[string]WinLossRecord {
    
    battleResults := make(chan WinLossRecord, 100)
    battlesToFight := make(chan []greedy.Strategy, 100)
    
    go func(){
        for i := 0; i < battlesPerGeneration; i++ {
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

func (rec WinLossRecord) Victory() float64 {
    return float64(rec.wins) / float64(rec.matches)
}

type SimpleStrategy struct {
    id string
    MinKeepPoints [8]int
    MinStopPoints [8]int
}

func (s SimpleStrategy) Id()  string {
    return s.id
}


func (s SimpleStrategy) Print() {
    fmt.Printf("Strategy for %s\n", s.id)
    fmt.Printf("Borrow with 1 die: %d; 2: %d; 3: %d;4: %d;5: %d;6: %d;7: %d;8: %d\n",
        s.MinKeepPoints[0], s.MinKeepPoints[1], s.MinKeepPoints[2], s.MinKeepPoints[3],
        s.MinKeepPoints[4], s.MinKeepPoints[5], s.MinKeepPoints[6], s.MinKeepPoints[7])

    fmt.Printf("Stop with 1 die: %d; 2: %d; 3: %d;4: %d;5: %d;6: %d;7: %d;8: %d\n",
        s.MinStopPoints[0], s.MinStopPoints[1], s.MinStopPoints[2], s.MinStopPoints[3],
        s.MinStopPoints[4], s.MinStopPoints[5], s.MinStopPoints[6], s.MinStopPoints[7])
}

func (s SimpleStrategy) ShouldKeep(dice greedy.Dice, game *greedy.Game, minScore int)  bool {
    return dice.RunningScore > s.MinKeepPoints[dice.Count - 1]
}

func (s SimpleStrategy) ShouldRoll(dice greedy.Dice, game *greedy.Game)  bool {
    return dice.RunningScore < s.MinStopPoints[dice.Count - 1]
}

func (a SimpleStrategy)  Breed(c greedy.Strategy, id string) greedy.Strategy {
    b, ok := c.(SimpleStrategy)
    if !ok {return nil}
    return SimpleStrategy {
        id: id,
        MinKeepPoints: averageLists(a.MinKeepPoints, b.MinKeepPoints),
        MinStopPoints: averageLists(a.MinStopPoints, b.MinStopPoints),
    }
}

func averageLists(a, b [8]int) [8]int {
    c := [8]int{}
    for i:=0; i<8;i++ {
        if rand.Intn(2) == 1 {
            c[i] = a[i]
        } else {
            c[i] = b[i]
        }
    }
    return c
}

func averageInts(a, b int) int {
    return (a + b) / 2
}

func SimpleStrategyGenerator(prefix string, count int) []greedy.Strategy {
    strategies := []greedy.Strategy{}
    for i := 0; i < count; i++ {
        strategies = append(strategies, greedy.Strategy(SimpleStrategy{
                prefix+":"+strconv.Itoa(i+1),
                [8]int{H(500,4000), H(500,4000), H(100,4000), H(100,4000),
                 H(100,4000),H(100,4000),H(100,3000),1},
                [8]int{H(0,3000), H(0,4000), H(0,5000), H(0,5000),
                H(500,6000),H(1000,8000),H(1000,8000),100000},
            }))
    }
    return strategies
}
func H(min, max int) int {
    return (rand.Intn((max - min) / 50) * 50) + min
}

func Helpful() int {
    return (rand.Intn(200) * 50)
}