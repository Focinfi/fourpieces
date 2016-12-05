# Four Pieces Chess

A Bot play chess with itself to improve the ability to win.

## Game Rule:

Board 4 x 4

![init](http://7xj8s4.com1.z0.glb.clouddn.com/four_pieces_init.jpg)

2-1 eat one piece

![](http://7xj8s4.com1.z0.glb.clouddn.com/four_pieces_w_1.jpg) 

![](http://7xj8s4.com1.z0.glb.clouddn.com/four_pieces_w_2.jpg) 

![](http://7xj8s4.com1.z0.glb.clouddn.com/four_pieces_w_3.jpg) 

2-2 eat none

![](http://7xj8s4.com1.z0.glb.clouddn.com/four_pieces_n_1.jpg) 
![](http://7xj8s4.com1.z0.glb.clouddn.com/four_pieces_n_2.jpg) 
![](http://7xj8s4.com1.z0.glb.clouddn.com/four_pieces_n_3.jpg) 
![](http://7xj8s4.com1.z0.glb.clouddn.com/four_pieces_n_4.jpg) 
![](http://7xj8s4.com1.z0.glb.clouddn.com/four_pieces_n_5.jpg) 
![](http://7xj8s4.com1.z0.glb.clouddn.com/four_pieces_n_6.jpg) 

## Project design

1. Save every step with the score into database.
2. Get all available steps.
3. Search the record with biggest probability of winning.
4. If the step leads to loss piece will lower socres.

## Have a try

```shell
$ go install github.com/Focinfi/fourpieces
$ cd $GOPATH/src/github.com/Focinfi/fourpieces
$ go test -run Play
```

When you first to run, the score will always be 0, cause our bot is just a fresher,
and every boy can make a fool decision, as following, move (0, 3) to (0, 2), and the 
A will eat the B(0, 2).

But that's ok, our bot B will take a note to avoid this step next time.

```shell
[FourPieces] otpSteps: (0, 3) => (0, 2) score: 0
[FourPieces] otpSteps: (1, 3) => (1, 2) score: 0
[FourPieces] otpSteps: (2, 3) => (3, 3) score: 0
[FourPieces] otpSteps: (3, 1) => (2, 1) score: 0
[FourPieces] otpSteps: (3, 1) => (3, 2) score: 0
[FourPieces] B: (0, 3) => (0, 2)
[FourPieces] >>>>>> T-2 <<<<<<
               [ -  A  B  -]
               [ A  -  -  B]
               [ -  -  A  B]
               [ A  B  -  -]
```

After you run many times of `go test -run Play`, our bot B is no longer a naive fresher, 
as following, it now is pretty experienced. 

```shell
[FourPieces] otpSteps: (0, 3) => (0, 2) score: -11
[FourPieces] otpSteps: (1, 3) => (1, 2) score: 0
[FourPieces] otpSteps: (2, 3) => (3, 3) score: 333
[FourPieces] otpSteps: (3, 1) => (2, 1) score: -11
[FourPieces] otpSteps: (3, 1) => (3, 2) score: 2
[FourPieces] B: (2, 3) => (3, 3)
[FourPieces] >>>>>> T-2 <<<<<<
               [ -  A  -  B]
               [ A  -  -  B]
               [ -  -  A  -]
               [ A  B  -  B]
```

1. (0, 3) => (0, 2) score: -11, B tried before and be warned this is a bad step.
2. (1, 3) => (1, 2) score: 0, obviously a bad step, but 0 means B haven't tried.
3. (2, 3) => (3, 3) score: 333, experience tells us this is a good choice.
4. (3, 1) => (2, 1) score: -11, just like 1.
5. (3, 1) => (3, 2) score: 2, seems normal.

It turns that our bot can make a more reasonable choice as experience accumulated.

## Need more experience analysis
Currently, bot only remember one step feedback.