Az Óbudai Egyetem Mérnökinforamtikus hallgatói vagyunk és a beadandó projektünk git repójátban barangolsz éppen :).
# Mi is ez?
A tervünk az, hogy ezt a webservert expose-oljuk és a címből egy qr-kódot generálunk amelyet kiragasztunk kampusz szerte, és meglátjuk, hogy ki fogja beolvasni

# Alap ötlet
Pár éve történt, hogy amerikai heckerek, a parkoló melleti oszlopokra kiragaszttak qr-kódokat, azzal a szöveggel, hogy "Itt tudod megkönnyíteni a fizeteséd!", nyilván angolul ;).
Ezen ötlet alapján kezdtünk mink is neki a feladatnak. 

Ezzel egyideűleg, egy kérdőívet is elindítunk, hogy felmérjük, hogy mi az emberek elméleti tudása a témában, és ezzel szembe állítjuk, hogy hányan fognak bele esni a mi csapdánkba.


# Követelmény
El is felejtettem leírni, hogy hogy tudja a nagyérdemű is kipróbálni a kis projektet.
Én a **go.1.24.1**-es verziót használtam a kód írásához. 

Ez szükséges a futtatáshoz is!

```bash
git clone https://github.com/ProfesszorPeter/webserver.git
cd webserver
go run main.go
```

Amennyiben az elveteműlt olvasó szeretne futtatható állomány készíteni:
```bash
git clone https://github.com/ProfesszorPeter/webserver.git
cd webserver
go build main.go
./main
```
## TODO
A későbbiekben lehet feltöltöm majd a futtatható állományt is, de nem gondolom, hogy az szükséges lenne.

