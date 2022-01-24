# Go SSE

Server Sent Events backend written in go. Redis is used to distribute published messages across all backends so that the backend can be scaled horizontally.

## Docker Compose

Start the compose stack, then visit the webpage at http://localhost:8080. It is a basic chat page. Open it multiple times to see real time events published from one client on all other client pages.

```bash
docker-compose up
```

## Example Output

```bash
backend_2  | client connected: 192.168.16.2:37712
backend_1  | client connected: 192.168.16.2:47996
backend_2  | 2022/01/23 23:01:44 publishing message: {"user":"John","message":"Hi"}
backend_2  | 2022/01/23 23:01:44 receiving message: {"user":"John","message":"Hi"}
backend_1  | 2022/01/23 23:01:44 receiving message: {"user":"John","message":"Hi"}
backend_2  | 2022/01/23 23:01:53 publishing message: {"user":"Maria","message":"Hello :)"}
backend_1  | 2022/01/23 23:01:53 receiving message: {"user":"Maria","message":"Hello :)"}
backend_2  | 2022/01/23 23:01:53 receiving message: {"user":"Maria","message":"Hello :)"}
```
