APP_NAME:= SeedBot
IMAGE_NAME:= seed-bot
PLATFORMS:= windows linux darwin
SOURCE_DIR:= build

build_app:
	@echo "Building App $(APP_NAME)..."
	go build -o $(SOURCE_DIR)/$(APP_NAME).exe
	@echo "Building App $(APP_NAME) Successfully..."

clean:
	@echo "Cleaning App $(APP_NAME) From $(SOURCE_DIR)"; \
	rm -f $(SOURCE_DIR)/$(APP_NAME).exe
	@echo "Cleaning App $(APP_NAME) From $(SOURCE_DIR) Successfully..."

rebuild: clean build_app
run: rebuild
	@echo "Running App..."
	./$(SOURCE_DIR)/$(APP_NAME).exe

image:
	@echo "Building Docker Image $(APP_NAME)..."
	docker build -t $(IMAGE_NAME) .
	@echo "Building Docker Image $(APP_NAME) Successfully..."

remove:
	@echo "Removing Docker Image $(APP_NAME)..."
	docker rmi $(IMAGE_NAME) --force
	@echo "Removing Docker Image $(APP_NAME) Successfully..."

up:
	@echo "Run Container $(APP_NAME)..."
	docker-compose up -d

down:
	@echo "Down Container $(APP_NAME)..."
	docker-compose down
