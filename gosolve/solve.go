package main

import (
    "fmt"
    "strings"
    "math/rand"
    "time"
)

const (
    EmptyColour = iota
    White
    Black
    MaxColour
)

const (
    EmptyPiece = iota
    Knight
    Bishop
    Rook
    Queen
    King
    MaxPiece
)

type Colour int
type Piece int
type Bitboard int16

type ColouredPiece struct {
    colour Colour
    piece Piece
}

type Puzzle struct {
    whitepieces Bitboard
    blackpieces Bitboard
    position [16]ColouredPiece
    whitekingsquare uint
    turn Colour
    turns int
}

var emptysquare = ColouredPiece{EmptyColour, EmptyPiece}
var whiteking = ColouredPiece{White, King}
var _movedests [16][MaxPiece]Bitboard
var _squaresbetween [16][16]Bitboard
var other [3]Colour = [3]Colour{White: Black, Black: White}

func main() {
    init_movedests()
    init_squaresbetween()

    solved := 0
    played := 0

    start := time.Now()

    for {
        puzzle := init_puzzle()

        solvable := puzzle.dfs()
        if solvable {
            solved++
        } else {
            fmt.Println("Unsolvable: ", puzzle.fen())
        }

        played++

        if played % 100000 == 0 {
            elapsed := time.Since(start)
            fmt.Println("Solved ", solved, " out of ", played, " in ", elapsed)
        }
    }
}

func (b Bitboard) str() string {
    s := ""
    for i := uint(0); i < 16; i++ {
        if i % 4 == 0 && i != 0{
            s += "/"
        }

        if (b & (1 << i)) != 0 {
            s += "1"
        } else {
            s += "."
        }
    }
    return s
}

func init_puzzle() Puzzle {
    var p Puzzle

    p.position = [16]ColouredPiece{
     ColouredPiece{White, Knight},
     ColouredPiece{White, Knight},
     ColouredPiece{Black, Knight},
     ColouredPiece{Black, Knight},
     ColouredPiece{White, Bishop},
     ColouredPiece{White, Bishop},
     ColouredPiece{Black, Bishop},
     ColouredPiece{Black, Bishop},
     ColouredPiece{White, Rook},
     ColouredPiece{White, Rook},
     ColouredPiece{Black, Rook},
     ColouredPiece{Black, Rook},
     ColouredPiece{White, Queen},
     ColouredPiece{Black, Queen},
     ColouredPiece{White, King},
     ColouredPiece{Black, King},
    }

    rand.Shuffle(len(p.position), func(i, j int) {
        p.position[i],p.position[j] = p.position[j],p.position[i]
    })

    p.turn = White
    p.turns = 0
    p.whitepieces = 0
    p.blackpieces = 0
    for i := uint(0); i < 16; i++ {
        if p.position[i].colour == White {
            p.whitepieces |= 1 << i
            if p.position[i].piece == King {
                p.whitekingsquare = i
            }
        } else {
            p.blackpieces |= 1 << i
        }
    }

    return p
}

func (p *Puzzle) fen() string {
    m := make(map[Piece]string)
    m[King] = "k"
    m[Queen] = "q"
    m[Rook] = "r"
    m[Bishop] = "b"
    m[Knight] = "n"
    m[EmptyPiece] = "1"

    fen := ""

    for i := uint(0); i < 16; i++ {
        c := m[p.position[i].piece]
        if p.position[i].colour == White {
            c = strings.ToUpper(c)
        }

        if i % 4 == 0 && i != 0{
            fen += "/"
        }
        fen += c
    }

    return fen
}

func (p *Puzzle) dfs() bool {
    //fmt.Println(p.fen())

    // if the king is the only remaining piece, we've solved it
    if ((p.whitepieces ^ (1 << p.whitekingsquare)) == 0) && p.blackpieces == 0 {
        //fmt.Println("SOLVED")
        return true
    }

    // if no pieces are adjacent to the white king, it's not solvable
    if (movedests(p.whitekingsquare, King) & (p.whitepieces | p.blackpieces)) == 0 {
        //fmt.Println("NOPE")
        return false
    }

    var fromboard *Bitboard
    var toboard *Bitboard

    if p.turn == White {
        fromboard = &p.whitepieces
        toboard = &p.blackpieces
    } else {
        fromboard = &p.blackpieces
        toboard = &p.whitepieces
    }

    nmoves := 0

    for from := uint(0); from < 16; from++ {
        if (*fromboard & (1 << from)) == 0 {
            continue;
        }
        //fmt.Println("movedests from ", from, " on ", p.fen(), " are ", movedests(from, p.position[from].piece).str())
        viablemoves := movedests(from, p.position[from].piece) & *toboard & ^(1 << p.whitekingsquare)
        //fmt.Println("Viable moves from ", from, " on ", p.fen(), " are ", viablemoves.str())
        for to := uint(0); to < 16; to++ {
            if (viablemoves & (1 << to)) == 0 {
                continue;
            }

            // can't jump over any other pieces
            if (squaresbetween(from, to) & (p.whitepieces | p.blackpieces)) != 0 {
                //fmt.Println("Can't move from ", from, " to ", to, " because there are pieces in the way on ", p.fen())
                continue;
            }

            nmoves++

            *fromboard ^= 1 << from
            *fromboard ^= 1 << to
            *toboard ^= 1 << to
            towas := p.position[to]
            p.position[to] = p.position[from]
            p.position[from] = emptysquare
            p.turn = other[p.turn]
            p.turns++

            if p.position[to] == whiteking {
                p.whitekingsquare = to
            }

            //fmt.Println("BEFORE: ", p.fen())
            solvable := p.dfs()
            //fmt.Println(" AFTER: ", p.fen())

            p.position[from] = p.position[to]
            if p.position[from] == whiteking {
                p.whitekingsquare = from
            }
            p.position[to] = towas
            *toboard ^= 1 << to
            *fromboard ^= 1 << to
            *fromboard ^= 1 << from
            p.turn = other[p.turn]
            p.turns--

            if solvable {
                //fmt.Println("TRUE")
                return true
            }
        }
    }

    //fmt.Println("RETURN AFTER ", nmoves, " MOVES ON ", p.fen())

    return false
}

func init_squaresbetween() {
    // XXX: should generate these automatically

    // horizontal rook moves
    _squaresbetween[0][2] = 1 << 1
    _squaresbetween[0][3] = (1 << 1) | (1 << 2)
    _squaresbetween[1][3] = 1 << 2
    _squaresbetween[4][6] = 1 << 5
    _squaresbetween[4][7] = (1 << 5) | (1 << 6)
    _squaresbetween[5][7] = 1 << 6
    _squaresbetween[8][10] = 1 << 9
    _squaresbetween[8][11] = (1 << 9) | (1 << 10)
    _squaresbetween[9][11] = 1 << 10
    _squaresbetween[12][14] = 1 << 13
    _squaresbetween[12][15] = (1 << 13) | (1 << 14)
    _squaresbetween[13][15] = 1 << 14

    // vertical rook moves
    _squaresbetween[0][8] = 1 << 4
    _squaresbetween[0][12] = (1 << 4) | (1 << 8)
    _squaresbetween[4][12] = 1 << 8
    _squaresbetween[1][9] = 1 << 5
    _squaresbetween[1][13] = (1 << 5) | (1 << 9)
    _squaresbetween[5][13] = 1 << 9
    _squaresbetween[2][10] = 1 << 6
    _squaresbetween[2][14] = (1 << 6) | (1 << 10)
    _squaresbetween[6][14] = 1 << 10
    _squaresbetween[3][11] = 1 << 7
    _squaresbetween[3][15] = (1 << 7) | (1 << 11)
    _squaresbetween[7][15] = 1 << 11

    // up-right bishop moves
    _squaresbetween[8][2] = 1 << 5
    _squaresbetween[12][6] = 1 << 9
    _squaresbetween[12][3] = (1 << 9) | (1 << 6)
    _squaresbetween[9][3] = 1 << 6
    _squaresbetween[13][7] = 1 << 10

    // down-right bishop moves
    _squaresbetween[1][11] = 1 << 6
    _squaresbetween[0][10] = 1 << 5
    _squaresbetween[0][15] = (1 << 5) | (1 << 10)
    _squaresbetween[5][15] = 1 << 10
    _squaresbetween[4][14] = 1 << 9

    // add the inverses
    for from := 0; from < 16; from++ {
        for to := 0; to < 16; to++ {
            if _squaresbetween[from][to] != 0 {
                _squaresbetween[to][from] = _squaresbetween[from][to]
            }
        }
    }
}

// return the set of squares between from and to, not including the endpoints
// (works for diagonal and straight lines only)
func squaresbetween(from uint, to uint) Bitboard {
    return _squaresbetween[from][to]
}

func init_movedests() {
    for from := uint(0); from < 16; from++ {
        fromx := from % 4
        fromy := from / 4

        for to := uint(0); to < 16; to++ {
            if to == from {
                continue
            }
            tox := to % 4
            toy := to / 4

            diffx := int(tox) - int(fromx)
            diffy := int(toy) - int(fromy)
            if diffx < 0 {
                diffx = -diffx
            }
            if diffy < 0 {
                diffy = -diffy
            }

            if diffx <= 1 && diffy <= 1 {
                _movedests[from][King] |= 1 << to
            }
            if (diffx == 2 && diffy == 1) || (diffx == 1 && diffy == 2) {
                _movedests[from][Knight] |= 1 << to
            }
            if (diffx == diffy) {
                _movedests[from][Bishop] |= 1 << to
                _movedests[from][Queen] |= 1 << to
            }
            if (diffx == 0 || diffy == 0) {
                _movedests[from][Rook] |= 1 << to
                _movedests[from][Queen] |= 1 << to
            }
        }
    }
}

// return the set of squares that the given piece on the given square can move to,
// assuming an empty board
func movedests(square uint, piece Piece) Bitboard {
    return _movedests[square][piece]
}
