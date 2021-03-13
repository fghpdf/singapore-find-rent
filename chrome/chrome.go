package chrome

import (
	"context"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/corpix/uarand"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	log "github.com/sirupsen/logrus"
)

type Condo struct {
	Name     string `bson:"Name"`
	Address  string `bson:"address"`
	District string `bson:"district"`
	// 资产期限
	Tenure string `bson:"tenure"`
	// 开发商
	Developer string    `bson:"developer"`
	Url       string    `bson:"url"`
	Facility  *Facility `bson:"facility, flatten"`
	FacString string    `bson:"fac_string"`
}

type Facility struct {
	// 泳池
	Pool bool `bson:"Pool"`
	// 网球场
	TennisCourt bool `bson:"tennis_court"`
	// 读书角
	ReadingCorner bool `bson:"reading_court"`
	// 屋顶花园
	RooftopGarden bool `bson:"rooftop_garden"`
	// 健身区
	FitnessArea bool `bson:"fitness_area"`
	// 俱乐部
	ClubHouse bool `bson:"club_house"`
	// 健身房
	Gymnasium bool `bson:"gymnasium"`
	// 烧烤设备
	BbqPit bool `bson:"bbq_pit"`
	// 秘密花园
	SecretGarden bool `bson:"secret_garden"`
	// 慢跑道
	JoggingTrack bool `bson:"jogging_track"`
	// 蒸汽室
	SteamRoom bool `bson:"steam_room"`
	// 停车场
	CarPark bool `bson:"car_park"`
	// 安保
	Security bool `bson:"security"`
}

func Run() {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent(uarand.GetRandom()),
	)

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)

	// create chrome instance
	ctx, cancel = chromedp.NewContext(
		ctx,
		chromedp.WithLogf(log.Infof),
	)
	defer cancel()

	run(ctx)
}

func RunWithRemote() {
	debugUrl := getDebugURL()

	ctx, cancel := chromedp.NewRemoteAllocator(context.Background(), debugUrl)
	defer cancel()

	run(ctx)
}

func run(ctx context.Context) {
	// create a timeout
	ctx, cancel := context.WithTimeout(ctx, 300*time.Second)
	defer cancel()
	dirUrls := analyzeDirectory(ctx, "https://condo.singaporeexpats.com/%sname/c")
	if len(dirUrls) == 0 {
		log.Fatalln("no condo")
	}

	var wg sync.WaitGroup

	condoUrls := make([]string, 0)
	for _, dirUrl := range dirUrls {
		wg.Add(1)
		go analyzeCondoList(ctx, dirUrl, &wg, &condoUrls)
	}

	wg.Wait()

	log.WithField("condo number", len(condoUrls)).Info("need to analyze")

	condos := loopAnalyzeCondo(ctx, condoUrls, 10)

	wg.Wait()

	mongoCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client, err := mongo.Connect(mongoCtx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Panic(err)
	}
	collection := client.Database("condo").Collection("condos")

	for _, condo := range condos {
		log.Info(condo)
		err := condo.insert(mongoCtx, collection)
		if err != nil {
			log.Errorf("mongo %v", err)
			continue
		}
	}

}
