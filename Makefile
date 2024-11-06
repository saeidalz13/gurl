build_mac_arm:
	GOOS=darwin go build -o gurl ./cmd

put_app_bin:
	sudo mv gurl /usr/local/bin