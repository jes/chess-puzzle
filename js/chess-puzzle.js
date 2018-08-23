var board = ChessBoard('board', {
    draggable: true,
    onDrop: piece_moved,
    showNotation: false,
});

var other = { 'w': 'b', 'b': 'w' };
var turn;

function init_puzzle() {
    var pieces = "KkQqRRrrBBbbNNnn".split("");
    pieces = shuffle(pieces);
    var fen = '';
    for (var i = 0; i < pieces.length; i++) {
        if (i > 0 && i%4 == 0)
            fen += '/';
        fen += pieces[i];
    }
    board.position(fen);
    turn = 'w';
}

// https://stackoverflow.com/a/6274398
function shuffle(array) {
    let counter = array.length;

    // While there are elements in the array
    while (counter > 0) {
        // Pick a random index
        let index = Math.floor(Math.random() * counter);

        // Decrease counter by 1
        counter--;

        // And swap the last element with it
        let temp = array[counter];
        array[counter] = array[index];
        array[index] = temp;
    }

    return array;
}

function dfs(position) {
    if (completed(position)) {
        return {
            moves: [],
            solved: true,
        };
    }

    // get set of occupied tiles from current position
    var tiles = [];
    for (var tile in position) {
        tiles.push(tile);
    }

    // just mindlessly try to move every tile to every other tile
    for (var i = 0; i < tiles.length; i++) {
        var from = tiles[i];
        for (var j = 0; j < tiles.length; j++) {
            var to = tiles[j];
            if (legal_move(position, from, to)) {
                var move = from + '-' + to;

                // apply the move
                var towas = position[to];
                position[to] = position[from];
                delete position[from];
                turn = other[turn];

                var r = dfs(position);

                // undo the move
                position[from] = position[to];
                position[to] = towas;
                turn = other[turn];

                if (r.solved) {
                    r.moves.unshift(move);
                    return {
                        moves: r.moves,
                        solved: true,
                    };
                }
            }
        }
    }

    return {
        solved: false,
    };
}

function solve() {
    var r = dfs(board.position());
    if (!r.solved) {
        alert("Not solvable!");
        return;
    }

    window.setTimeout(function() {
        showNextMove(r.moves);
    }, 500);
}

function showNextMove(moves) {
    var m = moves.shift();
    board.move(m);
    window.setTimeout(function() {
        showNextMove(moves);
    }, 500);
}

function completed(pos) {
    var pieces = '';
    for (var i in pos) {
        pieces += pos[i];
    }

    if (pieces == 'wK')
        return true;
    else
        return false;
}

function piece_moved(from, to, piece, newpos, oldpos, orientation) {
    if (!legal_move(oldpos, from, to))
        return 'snapback';

    // other colour's turn now
    turn = other[turn];
}

function legal_move(position, from, to) {
    if (!(from in position))
        return false;

    var piece = position[from];
    colour = piece.charAt(0);
    piece = piece.charAt(1);

    if (colour != turn)
        return false;

    // every move must be a capture
    if (!(to in position))
        return false;

    // must capture a piece of opposite colour
    var captured = position[to];
    var capturedcolour = captured.charAt(0);
    var capturedpiece = captured.charAt(1);
    if (capturedcolour == colour)
        return false;

    // can't capture white king
    if (capturedcolour == 'w' && capturedpiece == 'K')
        return false;

    // now validate the actual move mechanics
    var fromfile = from.charCodeAt(0) - 'a'.charCodeAt(0);
    var fromrank = from.charCodeAt(1) - '1'.charCodeAt(0);
    var tofile = to.charCodeAt(0) - 'a'.charCodeAt(0);
    var torank = to.charCodeAt(1) - '1'.charCodeAt(0);

    var signedmovefile = tofile - fromfile;
    var signedmoverank = torank - fromrank;

    var movefile = Math.abs(signedmovefile);
    var moverank = Math.abs(signedmoverank);

    // work out if the piece jumped over any other pieces (only applicable
    // to straight and diagonal moves, not knight's moves)
    var jumpedover = false;
    if ((movefile > 1 || moverank > 1) && (movefile == moverank || movefile == 0 || moverank == 0)) {
        var fdir = signedmovefile > 0 ? 1 : signedmovefile < 0 ? -1 : 0;
        var rdir = signedmoverank > 0 ? 1 : signedmoverank < 0 ? -1 : 0;
        var f = fromfile + fdir;
        var r = fromrank + rdir;
        while (f != tofile || r != torank) {
            var tile = String.fromCharCode(f + 'a'.charCodeAt(0)) + String.fromCharCode(r + '1'.charCodeAt(0));
            if (tile in position)
                jumpedover = true;
            f += fdir;
            r += rdir;
        }
    }

    // king
    if (piece == 'K') {
        if (movefile > 1 || moverank > 1)
            return false;
    }
    // queen
    if (piece == 'Q') {
        if (jumpedover || (movefile != 0 && moverank != 0 && movefile != moverank))
            return false;
    }
    // rook
    if (piece == 'R') {
        if (jumpedover || (movefile != 0 && moverank != 0))
            return false;
    }
    // bishop
    if (piece == 'B') {
        if (jumpedover || movefile != moverank)
            return false;
    }
    // knight
    if (piece == 'N') {
        if (!((movefile == 2 && moverank == 1) || (movefile == 1 && moverank == 2)))
            return false;
    }

    return true;
}

$('#solve').click(function() {
    solve();
});

$('#reset').click(function() {
    init_puzzle();
});

init_puzzle();
