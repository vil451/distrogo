package image

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"math/rand"
	"strings"
	"time"
)

var adjectives = []string{"admiring", "adoring", "affectionate", "agitated", "amazing", "angry", "awesome", "blissful", "boring", "brave", "clever", "cool", "compassionate", "competent", "confident", "cranky", "crazy", "dazzling", "determined", "distracted", "dreamy", "eager", "ecstatic", "elastic", "elated", "elegant", "eloquent", "epic", "fervent", "festive", "flamboyant", "focused", "friendly", "frosty", "gallant", "gifted", "goofy", "gracious", "happy", "hardcore", "heuristic", "hopeful", "hungry", "infallible", "inspiring", "jolly", "jovial", "keen", "kind", "laughing", "loving", "lucid", "magical", "mystifying", "modest", "musing", "naughty", "nervous", "nifty", "nostalgic", "objective", "optimistic", "peaceful", "pedantic", "pensive", "practical", "priceless", "quirky", "quizzical", "recursing", "relaxed", "reverent", "romantic", "sad", "serene", "sharp", "silly", "sleepy", "stoic", "stupefied", "suspicious", "sweet", "tender", "thirsty", "trusting", "unruffled", "upbeat", "vibrant", "vigilant", "vigorous", "wizardly", "wonderful", "xenodochial", "youthful", "zealous", "zen"}

var animals = []string{"aardvark", "albatross", "alligator", "ant", "anteater", "antelope", "ape", "armadillo", "donkey", "baboon", "badger", "barracuda", "bat", "bear", "beaver", "bee", "bison", "boar", "buffalo", "butterfly", "camel", "capybara", "caribou", "cassowary", "cat", "caterpillar", "cattle", "chamois", "cheetah", "chicken", "chimpanzee", "chinchilla", "chough", "clam", "cobra", "cockroach", "cod", "cormorant", "coyote", "crab", "crane", "crocodile", "crow", "curlew", "deer", "dinosaur", "dog", "dogfish", "dolphin", "dotterel", "dove", "dragonfly", "duck", "dugong", "dunlin", "eagle", "echidna", "eel", "eland", "elephant", "elk", "emu", "falcon", "ferret", "finch", "fish", "flamingo", "fly", "fox", "frog", "gaur", "gazelle", "gerbil", "giant panda", "giraffe", "gnat", "gnu", "goat", "goldfinch", "goldfish", "goose", "gorilla", "goshawk", "grasshopper", "grouse", "guanaco", "guinea fowl", "guinea pig", "gull", "hamster", "hare", "hawk", "hedgehog", "heron", "herring", "hippopotamus", "hornet", "horse", "human", "hummingbird", "hyena", "ibex", "ibis", "jackal", "jaguar", "jay", "jellyfish", "kangaroo", "kingfisher", "koala", "komodo dragon", "kookabura", "kouprey", "kudu", "lapwing", "lark", "lemur", "leopard", "lion", "llama", "lobster", "locust", "loris", "louse", "lyrebird", "magpie", "mallard", "manatee", "mandrill", "mantis", "marten", "meerkat", "mink", "mole", "mongoose", "monkey", "moose", "mosquito", "mouse", "mule", "narwhal", "newt", "nightingale", "octopus", "okapi", "opossum", "oryx", "ostrich", "otter", "owl", "ox", "oyster", "panther", "parrot", "partridge", "peafowl", "pelican", "penguin", "pheasant", "pig", "pigeon", "polar bear", "pony", "porcupine", "porpoise", "quail", "quelea", "quetzal", "rabbit", "raccoon", "rail", "ram", "rat", "raven", "red deer", "red panda", "reindeer", "rhinoceros", "rook", "salamander", "salmon", "sand dollar", "sandpiper", "sardine", "scorpion", "sea lion", "sea urchin", "seahorse", "seal", "shark", "sheep", "shrew", "skunk", "snail", "snake", "sparrow", "spider", "spoonbill", "squid", "squirrel", "starling", "stingray", "stinkbug", "stork", "swallow", "swan", "tapir", "tarsier", "termite", "tiger", "toad", "trout", "turkey", "turtle", "viper", "vulture", "wallaby", "walrus", "wasp", "weasel", "whale", "wildcat", "wolf", "wolverine", "wombat", "woodcock", "woodpecker", "worm", "wren", "yak", "zebra"}

var defaultImage = "registry.fedoraproject.org/fedora-toolbox:39"

type Image struct {
	imageName     string
	containerName string
}

func New(imageName string, containerName string) (*Image, error) {
	if imageName == "" {
		imageName = defaultImage
	} else {
		imageName = ensureTag(imageName)
	}

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	images, err := cli.ImageList(ctx, types.ImageListOptions{})

	if err != nil {
		return nil, err
	}

	if containerName == "" {
		containerName = generateContainerName()
	}

	for _, image := range images {
		for _, tag := range image.RepoTags {
			if tag == imageName {
				return &Image{
					imageName:     imageName,
					containerName: containerName,
				}, nil
			}
		}
	}

	return nil, nil
}

func generateContainerName() string {
	rand.Seed(time.Now().UnixNano())
	adj := adjectives[rand.Intn(len(adjectives))]
	animal := animals[rand.Intn(len(animals))]
	return fmt.Sprintf("%s_%s", adj, animal)
}

func hasTag(imageName string) bool {
	return strings.Contains(imageName, ":")
}

func ensureTag(imageName string) string {
	if !hasTag(imageName) {
		return imageName + ":latest"
	}
	return imageName
}
