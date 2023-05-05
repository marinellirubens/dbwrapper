curl -s -w '%{time_total}' http://localhost:8080/pg\?query\="select+*+from+public.Employee" -o /dev/null
