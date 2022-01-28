Тестовое задание для FindMyKids 

Как использовать: 

    -установить зависимости 
    -запустить parent/notify.go, он будет слушать приходящие координаты и определять, вышел ли ребенок из радиуса
    -запутить child/child.go для теста

**Можно скомпилировать через Makefile, или запустить из консоли. Ниже инструкция для двух вариантов**

Запуск:
    
    -(notify_parent_simulator_darwin_amd64 или go run cmd/parent/notify.go) start -addr 127.0.0.1:4079
    -(child_parent_simulator_darwin_amd64 или go run cmd/child/child.go) start -target http://127.0.0.1:4079/set-coord -addr 127.0.0.1:4080