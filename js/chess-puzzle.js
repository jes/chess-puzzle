var board = ChessBoard('board', {
    draggable: true,
    onDrop: piece_moved,
});

var other = { 'w': 'b', 'b': 'w' };
var turn;

function init_puzzle() {
    board.position("kNbr/KnBR/qNbr/QnBR");
    turn = 'w';
}

function piece_moved(from, to, piece, newpos, oldpos, orientation) {
    if (!legal_move(oldpos, from, to))
        return 'snapback';

    // other colour's turn now
    turn = other[turn];

    // TODO: detect end of game (either white king is the only piece, or there are no legal
    // moves)
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

init_puzzle();
