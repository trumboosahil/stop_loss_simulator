# Stop Loss Simulator

## Overview

The Stop Loss Simulator is a Go application designed to simulate order and tick events for financial instruments. It uses Redis to manage and process stop-loss orders efficiently.

## Features

- Preloads millions of orders with stop-loss conditions into Redis.
- Simulates tick events for predefined symbols.
- Publishes tick events to Redis Pub/Sub.

## Prerequisites

- Go 1.18 or later
- Redis server running locally or accessible remotely
- Git for version control

## Installation

1. **Clone the repository:**

   ```bash
   git clone git@github-personal:trumboosahil/stop_loss_simulator.git
   cd stop_loss_simulator
   ```
2. **Install dependencies:**

   Ensure you have Go installed and set up properly. You can install dependencies using:

   ```bash
   go mod tidy
   ```
3. **Run the application:**

   ```bash
   go run main.go
   ```

## Usage

- **Preload Orders:** The application will automatically preload orders into Redis when started.
- **Simulate Tick Events:** Tick events are simulated at a regular interval and published to Redis.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request for any improvements or bug fixes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contact

For any questions or issues, please contact [trumboosahil](mailto:your-email@example.com).
