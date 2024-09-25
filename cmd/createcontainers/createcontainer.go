package createcontainers

import (
	"context"
	"distrogo/cmd/dockerclient"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

var adjectives = []string{"admiring", "adoring", "affectionate", "agitated", "amazing", "angry", "awesome", "blissful", "boring", "brave", "clever", "cool", "compassionate", "competent", "confident", "cranky", "crazy", "dazzling", "determined", "distracted", "dreamy", "eager", "ecstatic", "elastic", "elated", "elegant", "eloquent", "epic", "fervent", "festive", "flamboyant", "focused", "friendly", "frosty", "gallant", "gifted", "goofy", "gracious", "happy", "hardcore", "heuristic", "hopeful", "hungry", "infallible", "inspiring", "jolly", "jovial", "keen", "kind", "laughing", "loving", "lucid", "magical", "mystifying", "modest", "musing", "naughty", "nervous", "nifty", "nostalgic", "objective", "optimistic", "peaceful", "pedantic", "pensive", "practical", "priceless", "quirky", "quizzical", "recursing", "relaxed", "reverent", "romantic", "sad", "serene", "sharp", "silly", "sleepy", "stoic", "stupefied", "suspicious", "sweet", "tender", "thirsty", "trusting", "unruffled", "upbeat", "vibrant", "vigilant", "vigorous", "wizardly", "wonderful", "xenodochial", "youthful", "zealous", "zen"}

var animals = []string{"aardvark", "albatross", "alligator", "ant", "anteater", "antelope", "ape", "armadillo", "donkey", "baboon", "badger", "barracuda", "bat", "bear", "beaver", "bee", "bison", "boar", "buffalo", "butterfly", "camel", "capybara", "caribou", "cassowary", "cat", "caterpillar", "cattle", "chamois", "cheetah", "chicken", "chimpanzee", "chinchilla", "chough", "clam", "cobra", "cockroach", "cod", "cormorant", "coyote", "crab", "crane", "crocodile", "crow", "curlew", "deer", "dinosaur", "dog", "dogfish", "dolphin", "dotterel", "dove", "dragonfly", "duck", "dugong", "dunlin", "eagle", "echidna", "eel", "eland", "elephant", "elk", "emu", "falcon", "ferret", "finch", "fish", "flamingo", "fly", "fox", "frog", "gaur", "gazelle", "gerbil", "giant panda", "giraffe", "gnat", "gnu", "goat", "goldfinch", "goldfish", "goose", "gorilla", "goshawk", "grasshopper", "grouse", "guanaco", "guinea fowl", "guinea pig", "gull", "hamster", "hare", "hawk", "hedgehog", "heron", "herring", "hippopotamus", "hornet", "horse", "human", "hummingbird", "hyena", "ibex", "ibis", "jackal", "jaguar", "jay", "jellyfish", "kangaroo", "kingfisher", "koala", "komodo dragon", "kookabura", "kouprey", "kudu", "lapwing", "lark", "lemur", "leopard", "lion", "llama", "lobster", "locust", "loris", "louse", "lyrebird", "magpie", "mallard", "manatee", "mandrill", "mantis", "marten", "meerkat", "mink", "mole", "mongoose", "monkey", "moose", "mosquito", "mouse", "mule", "narwhal", "newt", "nightingale", "octopus", "okapi", "opossum", "oryx", "ostrich", "otter", "owl", "ox", "oyster", "panther", "parrot", "partridge", "peafowl", "pelican", "penguin", "pheasant", "pig", "pigeon", "polar bear", "pony", "porcupine", "porpoise", "quail", "quelea", "quetzal", "rabbit", "raccoon", "rail", "ram", "rat", "raven", "red deer", "red panda", "reindeer", "rhinoceros", "rook", "salamander", "salmon", "sand dollar", "sandpiper", "sardine", "scorpion", "sea lion", "sea urchin", "seahorse", "seal", "shark", "sheep", "shrew", "skunk", "snail", "snake", "sparrow", "spider", "spoonbill", "squid", "squirrel", "starling", "stingray", "stinkbug", "stork", "swallow", "swan", "tapir", "tarsier", "termite", "tiger", "toad", "trout", "turkey", "turtle", "viper", "vulture", "wallaby", "walrus", "wasp", "weasel", "whale", "wildcat", "wolf", "wolverine", "wombat", "woodcock", "woodpecker", "worm", "wren", "yak", "zebra"}

var default_image = "registry.fedoraproject.org/fedora-toolbox:39"

func CreateContainer() *cobra.Command {
	var containerName string
	var imageName string
	var pullImage bool
	command := &cobra.Command{
		Use:     "create",
		Short:   "Create a container",
		Aliases: []string{"c"},
		Run: func(cmd *cobra.Command, args []string) {
			if imageName == "" {
				imageName = default_image
			} else {
				imageName = ensureTag(imageName)
			}
			create(imageName, containerName, pullImage)
		},
	}

	command.Flags().StringVarP(
		&imageName,
		"image",
		"i",
		"",
		"image name of a container",
	)

	command.Flags().StringVarP(
		&containerName,
		"name",
		"n",
		"",
		"container name",
	)

	command.Flags().BoolVarP(
		&pullImage,
		"pull",
		"p",
		false,
		"pull image",
	)
	return command
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

func isImageAvailableLocally(ctx context.Context, cli *client.Client, imageName string) bool {

	images, err := cli.ImageList(ctx, types.ImageListOptions{})

	if err != nil {
		return false
	}

	for _, image := range images {
		for _, tag := range image.RepoTags {
			if tag == imageName {
				return true
			}
		}
	}
	return false
}

func create(image string, containerName string, pull bool) {
	ctx := context.Background()
	cli, err := dockerclient.InitDockerClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating Docker client: %v\n", err)
		os.Exit(1)
	}
	defer dockerclient.CloseDockerClient(cli)

	if containerName == "" {
		containerName = generateContainerName()
	}

	if pull {
		_, err := pullImage(ctx, cli, image)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error pulling image: %v\n", err)
			os.Exit(1)
		}
	} else if !isImageAvailableLocally(ctx, cli, image) {
		fmt.Fprintf(os.Stderr, "Image not found, add --pull options for pulling!\n")
		return
	}

	_, err = createContainer(ctx, cli, image, containerName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating container: %v\n", err)
		os.Exit(1)
	}

	_, err = runContainer(ctx, cli, containerName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running container: %v\n", err)
		os.Exit(1)
	}
}

func runContainer(ctx context.Context, cli *client.Client, name string) (interface{}, error) {
	if err := cli.ContainerStart(ctx, name, container.StartOptions{}); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting container: %v\n", err)
		return nil, err
	}
	return cli.ContainerInspect(ctx, name)
}

func pullImage(ctx context.Context, cli *client.Client, name string) (io.ReadCloser, error) {
	config := &container.Config{
		Image: name,
	}
	resp, err := cli.ImagePull(ctx, config.Image, types.ImagePullOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error pulling image: %v\n", err)
		return nil, err
	}
	defer resp.Close()

	_, err = io.Copy(os.Stdout, resp)
	if err != nil {
		return nil, fmt.Errorf("error copying response: %v", err)
	}

	return resp, nil
}

func createContainer(ctx context.Context, cli *client.Client, image string, name string) (container.CreateResponse, error) {
	//options := types.ContainerListOptions{}
	config := &container.Config{
		Image: image,
		Cmd:   []string{"echo", "Hello, World!"},
	}
	hostConfig := &container.HostConfig{}
	resp, err := cli.ContainerCreate(ctx, config, hostConfig, nil, nil, name)
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Error creating container: %v\n", err)
		if err != nil {
			return container.CreateResponse{}, err
		}
		os.Exit(1)
	}

	return resp, nil
}
