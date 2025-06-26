+-------------------+
|     Manager       |
|-------------------|
| - clients         |<-------------------+
| - handlers        |                    |
|-------------------|                    |
| + addClient()     |                    |
| + removeClient()  |                    |
| + serveWS()       |                    |
| + routeEvent()    |                    |
+-------------------+                    |
         ^                               |
         |                               |
         | manager                       |
+-------------------+                    |
|     Client        |                    |
|-------------------|                    |
| - connection      |                    |
| - manager ------- |--------------------+
| - egress          |
|-------------------|
| + readMessages()  |<--- goroutine
| + writeMessages() |<--- goroutine
+-------------------+

[HTTP Request] ---> serveWS() ---> NewClient()
                                 addClient()
                                 go readMessages()
                                 go writeMessages()