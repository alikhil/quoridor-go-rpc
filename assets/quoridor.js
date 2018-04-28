//================================Models=======================================================================
style = {
    cellBaseColor: 0xbababa,
    cellLineColor: 0x333333,
    cellHoverColor: 0x848484,
    cellHighlightColor: 0x8de29,
    cellOverHighlighted: 0x0d9b23,

    wallTransparent: 0xededed,
    wallBaseColor: 0xa56445,
    wallLineColor: 0x333333,

    pawnColor: [0x50a545, 0x45a0a5, 0x7045a5, 0xa5454b],

    backgroundColor: 0xededed,

    windowWidth: 600,
    windowHeight: 600,
    cellWidth: 50,
    cellGap: 15
}


class Cell extends PIXI.Graphics {
    constructor(x, y, w, grid_pos) {
        super();

        this.x = x;
        this.y = y;
        this.w = w;

        this.pos = grid_pos;
        this.highlighted = false;
        this.myself = this;

        this.drawNormal();

        this.interactive = true;
        this.on('mouseout', this.onOut);
        this.on('mouseover', this.onOver);
        this.on('mousedown', this.onClick);

        app.stage.addChild(this);
    }

    onOver() {
        if (this.highlighted) {
            this.drawOverHighlighted()
        }
    }

    onOut() {
        if (this.highlighted) {
            this.drawHighlighted();
        }
    }

    onClick() {
        if (game.stage == 'move') {
            if (this.highlighted) {
                move_MovePawn(this.myself);
            } else {
                alert('You can\'t move into this cell');
            }
        } else {
            alert('You can\'t move now');
        }
    }

    highlight() {
        this.highlighted = true;
        this.drawHighlighted();
    }

    unhighlight() {
        this.highlighted = false;
        this.drawNormal();
    }

    drawNormal() {
        var w = this.w;
        this.clear()
        this.beginFill(style.cellBaseColor);
        this.drawRect(0, 0, w, w);
        this.endFill();
    }
    drawHighlighted() {
        var w = this.w;
        this.clear()
        this.lineStyle(1, style.cellLineColor, 1);
        this.beginFill(style.cellHighlightColor);
        this.drawRect(0, 0, w, w);
        this.endFill();
    }
    drawOverHighlighted() {
        var w = this.w;
        this.clear()
        this.lineStyle(1, style.cellLineColor, 1);
        this.beginFill(style.cellOverHighlighted);
        this.drawRect(0, 0, w, w);
        this.endFill();
    }

}

class Pawn extends PIXI.Graphics {
    constructor(grid_pos, id, taret_column, target_row) {
        super();

        this.id = id;
        this.pos = grid_pos;

        this.r = style.cellWidth / 2;
        this.drawNormal();
        this.updCoord();
        this.taret_column = taret_column;
        this.target_row = target_row;

        app.stage.addChild(this);
    }

    moveTo(cell) {
        this.pos.col = cell.pos.col;
        this.pos.row = cell.pos.row;

        if ((this.pos.col == this.taret_column) || (this.pos.row == this.target_row)) {
            winner(this.id);
            game.stage = 'end';
        }

        this.updCoord();
    }

    updCoord() {
        this.r = style.cellWidth / 2;
        this.x = this.pos.col * (style.cellWidth + style.cellGap) + style.cellGap + this.r;
        this.y = this.pos.row * (style.cellWidth + style.cellGap) + style.cellGap + this.r;
        this.position.set(this.x, this.y);

    }

    drawNormal() {
        var w = this.w;
        this.lineStyle(1, style.cellLineColor, 1);
        this.beginFill(style.pawnColor[this.id]);
        this.drawCircle(0, 0, this.r);
        this.endFill();

    }
}

class Wall extends PIXI.Graphics {
    constructor(x, y, w, l, grid_pos) {
        super();

        this.x = x;
        this.y = y;
        this.w = w;
        this.l = l;

        this.placed = false;
        this.possible = true;
        this.pos = grid_pos;
        this.myself = this;

        this.drawOut();

        this.interactive = true;
        this.on('mouseout', this.onOut);
        this.on('mouseover', this.onOver);
        this.on('click', this.onClick);
        this.on('dclick', this.onDClick);

        app.stage.addChild(this);
    }

    onOver() {
        if (!this.placed && this.possible) {
            this.drawFinal();
        }
    }
    onOut() {
        if (!this.placed && this.possible) {
            this.drawOut();
        }
    }

    onClick() {
        if (game.stage == 'move') {
            if (walls_left[turnNumber] > 0) {
                if (isWallValid(this.pos) && !this.placed) {
                    this.drawFinal();
                    move_PlaceWall(this.myself);
                } else {
                    alert("Invalid wall placement");
                }
            } else {
                alert("You don't have walls left")
            }
        } else {
            alert('You can\'t place the wall now!');
        }
    }
    onDClick() {
    }

    drawFinal() {
        var w = this.w;
        var l = this.l;
        this.clear()
        this.lineStyle(1, style.wallLineColor, 1);
        this.beginFill(style.wallBaseColor);
        this.drawRect(0, 0, w, l);
        this.endFill();
    }

    drawOut() {
        var w = this.w;
        var l = this.l;
        this.clear()
        this.beginFill(style.wallTransparent, 0);
        this.drawRect(0, 0, w, l);
        this.endFill();
    }

}

class AdjMatrix {
    // cell id is col * 10 + row
    // matrix[from][to]
    constructor() {
        this.matrix = new Array(90);
        for (var i = 0; i < 90; i++) {
            this.matrix[i] = new Array(90);
            for (var j = 0; j < 90; j++) {
                this.matrix[i][j] = false;
            }
        }
    }
    addConnection(from, to) {
        var a = this.convertToNum(from);
        var b = this.convertToNum(to);
        this.matrix[a][b] = true;
        this.matrix[b][a] = true;
    }

    delConnection(from, to) {
        var a = this.convertToNum(from);
        var b = this.convertToNum(to);
        this.matrix[a][b] = false;
        this.matrix[b][a] = false;
    }

    moveIsValid(from, to) {
        if (!this._posIsValid(from) || !this._posIsValid(to)) {
            return false;
        }
        return this.matrix[this.convertToNum(from)][this.convertToNum(to)];
    }

    getJumpPos(from, over) {
        var goal_pos = null;
        if (from.col == over.col) { // vertical jump
            if (from.row > over.row) { // jump left 
                goal_pos = { col: from.col, row: from.row - 2 };
            } else { // jump right
                goal_pos = { col: from.col, row: from.row + 2 };
            }
        } else { // horizontal jump
            if (from.col > over.col) { // jump up
                goal_pos = { col: from.col - 2, row: from.row };
            } else { // jump down
                goal_pos = { col: from.col + 2, row: from.row };
            }
        }
        if (!this._posIsValid(goal_pos)) {
            return null;
        }

        if (this.moveIsValid(over, goal_pos)) {
            return goal_pos;
        }
        // goal_pos = {col: from.col -2 , row: from.row};
        // return {row:1, col:1}
        // return null
    }

    pathExists(pawn) {
        var target_cells = []

        if (pawn.target_row != null) {
            for (var i = 0; i < 9; i++) {
                target_cells.push(this.convertToNum(grid[i][pawn.target_row].pos));
            }
        } else {
            for (var i = 0; i < 9; i++) {
                target_cells.push(this.convertToNum(grid[pawn.taret_column][i].pos));
            }
        }

        var visited_cells = [];
        var neighbour_cells = [];

        neighbour_cells.push(this.convertToNum(pawn.pos));

        while (neighbour_cells.length > 0) {
            var cell = neighbour_cells.pop();
            if (target_cells.includes(cell)) {
                return true;
            }

            for (var i = 0; i < this.matrix.length; i++) {
                if (this.matrix[cell][i] == true) {
                    if (target_cells.includes(i)) {
                        return true;
                    } else {
                        if (!visited_cells.includes(i)) {
                            neighbour_cells.push(i);
                            visited_cells.push(i);
                        }

                    }
                }
            }
        }
        return false;
    }

    convertToNum(pos) {
        return pos.col * 10 + pos.row;
    }
    convertToPos(num) {
        return ({ col: Math.floor(num / 10), row: num % 10 })
    }
    _posIsValid(pos) {
        return pos.col >= 0 && pos.row >= 0 && pos.col <= 8 && pos.row <= 8;
    }
}

//=================================Functions=============================================================================

function isWallValid(pos) {
    check = [];

    if (pos.orient == 'vert') {
        check.push(wallsVert[pos.col][pos.row - 1]);
        check.push(wallsVert[pos.col][pos.row + 1]);
        check.push(wallsHor[pos.col][pos.row]);
    } else {
        if (pos.col > 1) {
            check.push(wallsHor[pos.col - 1][pos.row]);
        }
        if (pos.col < 7) {
            check.push(wallsHor[pos.col + 1][pos.row]);
        }
        check.push(wallsVert[pos.col][pos.row]);
    }

    // placement probles detection
    found_err = false;
    check.forEach(element => {
        if (element != null) {
            if (element.placed) {
                found_err = true;
            }
        }
    });
    if (found_err) {
        return false;
    }

    // check that all pawns can achieve their destination
    placeElementWall(pos);
    found_err = false;
    pawns.forEach(pawn => {
        if (!pathfinder.pathExists(pawn)) {
            found_err = true;
        }
    })
    removeElementWall(pos);
    if (found_err) {
        return false;
    }


    check.forEach(element => {
        if (element != null) {
            element.possible = false;
        }
    });
    return true;

}

function pathExists() {
    return true;
}

function showPossibleMoves() {
    var pawn = pawns[turnNumber];
    var pos = pawn.pos;
    var col = pos.col;
    var row = pos.row;

    var possible_coordinates = [];
    var nearby_pawns_pos = [];

    possible_coordinates.push([col, row - 1]);
    possible_coordinates.push([col, row + 1]);
    possible_coordinates.push([col - 1, row]);
    possible_coordinates.push([col + 1, row]);

    possible_coordinates.forEach(coord => {
        if (pathfinder.moveIsValid(pos, { col: coord[0], row: coord[1] })) {
            if (!cellIsBusy({ col: coord[0], row: coord[1] })) {
                highlighted.push(grid[coord[0]][coord[1]]);
            } else {
                nearby_pawns_pos.push({ col: coord[0], row: coord[1] });
            }
        }
    });

    nearby_pawns_pos.forEach(pos_over => {
        var j = pathfinder.getJumpPos(pos, pos_over);
        if (j != null) {
            if (!cellIsBusy(j)) {
                highlighted.push(grid[j.col][j.row]);
            }
        }
    });

    highlighted.forEach(c => {
        c.highlight()
    })
}

function cellIsBusy(pos) {
    var busy = false;
    pawns.forEach(p => {
        if (p.pos.col == pos.col && p.pos.row == pos.row) {
            busy = true;
        }
    })
    return busy;
}

function hideHighlighted() {
    highlighted.forEach(c => {
        c.unhighlight()
    })
    highlighted = [];
}

function placeElementWall(pos) {
    if (pos.orient == 'vert') {
        pathfinder.delConnection(grid[pos.col][pos.row].pos, grid[pos.col + 1][pos.row].pos);
        pathfinder.delConnection(grid[pos.col][pos.row + 1].pos, grid[pos.col + 1][pos.row + 1].pos);
    } else {
        pathfinder.delConnection(grid[pos.col][pos.row].pos, grid[pos.col][pos.row + 1].pos);
        pathfinder.delConnection(grid[pos.col + 1][pos.row].pos, grid[pos.col + 1][pos.row + 1].pos);
    }
}

function removeElementWall(pos) {
    if (pos.orient == 'vert') {
        pathfinder.addConnection(grid[pos.col][pos.row].pos, grid[pos.col + 1][pos.row].pos);
        pathfinder.addConnection(grid[pos.col][pos.row + 1].pos, grid[pos.col + 1][pos.row + 1].pos);
    } else {
        pathfinder.addConnection(grid[pos.col][pos.row].pos, grid[pos.col][pos.row + 1].pos);
        pathfinder.addConnection(grid[pos.col + 1][pos.row].pos, grid[pos.col + 1][pos.row + 1].pos);
    }
}

function move_PlaceWall(wall) {
    wall.placed = true;
    placeElementWall(wall.pos);
    hideHighlighted();
    walls_left[turnNumber] -= 1;
    step = {
        step: gstepId,
        data: "w#"+wall.pos.row+"#"+wall.pos.col+"#"+wall.pos.orient
    }
    game.stage = 'moveDone';
}

function move_MovePawn(cell) {
    hideHighlighted();
    pawns[turnNumber].moveTo(cell);
    step = {
        step: gstepId,
        data: "p#"+turnNumber+"#"+cell.pos.col+"#"+cell.pos.orient
    }
    game.stage = 'moveDone';
}

function winner(id) {
    alert("" + id + "'th player wins");
    location.reload();
}

//=================================Initialise=============================================================================

//Create a Pixi Application
let app = new PIXI.Application({
    width: style.windowWidth,
    height: style.windowHeight,
    backgroundColor: style.backgroundColor,
});

var grid = new Array(9);
var wallsVert = new Array(8);
var wallsHor = new Array(8);
var pathfinder = new AdjMatrix();
var pawns = [];
var highlighted = [];
var walls_left = [];

var turnNumber = 0;
var total_players = 0;
var game = null;
var gstepId = 0;
var step = {};

//=================================Start=============================================================================
class Game {
    constructor() {
        //Add the canvas that Pixi automatically created for you to the HTML document
        document.body.appendChild(app.view);
        this.s = 'stop';
        showPossibleMoves();
    }
    get stage() {
        return this.s;
    }
    set stage(value) {
        this.s == value;
        if (value == 'moveDone') {
            // if (turnNumber == total_players - 1) {
            //     turnNumber = 0;
            // } else {
            //     turnNumber += 1;
            // }
            // showPossibleMoves();
            socket.emit('share_step', step);
        } else if (value == 'end') {
            alert('Game ended');
        }
    }
}

class Client {
    constructor() {
        this.players_n = total_players;
        this.initBoard();
        this.initPawns();
        game = new Game();
    }

    initBoard() {
        // cells
        for (var i = 0; i < 9; i++) {
            grid[i] = new Array(9);
            for (var j = 0; j < 9; j++) {
                var x = i * (style.cellWidth + style.cellGap) + style.cellGap;
                var y = j * (style.cellWidth + style.cellGap) + style.cellGap;

                grid[i][j] = new Cell(x, y, style.cellWidth, { col: i, row: j });
            }
        }
        // add connections
        for (var i = 0; i < 8; i++) {
            for (var j = 0; j < 8; j++) {
                pathfinder.addConnection(grid[i][j].pos, grid[i + 1][j].pos);
                pathfinder.addConnection(grid[i][j].pos, grid[i][j + 1].pos);
            }
        }
        // last column connections
        for (var i = 0; i < 8; i++) {
            pathfinder.addConnection(grid[8][i].pos, grid[8][i + 1].pos);
        }
        // last row connections
        for (var i = 0; i < 8; i++) {
            pathfinder.addConnection(grid[i][8].pos, grid[i + 1][8].pos);
        }
        // vertical walls
        for (var i = 1; i < 9; i++) {
            wallsVert[i - 1] = new Array(8);
            for (var j = 0; j < 8; j++) {
                var x = i * (style.cellWidth + style.cellGap);
                var y = j * (style.cellWidth + style.cellGap) + style.cellGap;


                wallsVert[i - 1][j] = new Wall(x, y, style.cellGap, style.cellWidth * 2 + style.cellGap, { col: i - 1, row: j, orient: 'vert' });
            }
        }
        // horizontal walls
        for (var i = 0; i < 8; i++) {
            wallsHor[i] = new Array(8);
            for (var j = 1; j < 9; j++) {
                var x = i * (style.cellWidth + style.cellGap) + style.cellGap;
                var y = j * (style.cellWidth + style.cellGap);


                wallsHor[i][j - 1] = new Wall(x, y, style.cellWidth * 2 + style.cellGap, style.cellGap, { col: i, row: j - 1, orient: 'hor' });
            }
        }
    }

    initPawns() {
        if (this.players_n == 2) {
            pawns.push(new Pawn({ col: 4, row: 0 }, 0, null, 8));
            pawns.push(new Pawn({ col: 4, row: 8 }, 1, null, 0));
            walls_left = [10, 10];
        } else if (this.players_n == 3) {
            pawns.push(new Pawn({ col: 4, row: 0 }, 1, null, 8));
            pawns.push(new Pawn({ col: 8, row: 4 }, 2, 0, null));
            pawns.push(new Pawn({ col: 4, row: 8 }, 3, null, 0));
            walls_left = [7, 7, 7];
        } else if (this.players_n == 4) {
            pawns.push(new Pawn({ col: 0, row: 4 }, 0, 8, null));
            pawns.push(new Pawn({ col: 4, row: 0 }, 1, null, 8));
            pawns.push(new Pawn({ col: 8, row: 4 }, 2, 0, null));
            pawns.push(new Pawn({ col: 4, row: 8 }, 3, null, 0));
            walls_left = [5, 5, 5, 5];
        }
    }
}

socket = io();

function subscribe(){
    socket.on('show_endpoint', console.log);
    socket.on('make_step', makeStep);
    socket.on('apply_step', applyStep);
    socket.on('show_error', showError);
}

function makeStep(stepId, index){
    gstepId = stepId;
    turnNumber = index;
    showPossibleMoves();
}

function showError(s){
    console.log(s);
}

function applyStep(s){
    console.log();
}

function createNewGame(form) {
    if (document.getElementById('rb2').checked) {
        total_players = 2;
    } else if (document.getElementById('rb3').checked) {
        total_players = 3;
    } else if (document.getElementById('rb4').checked) {
        total_players = 4;
    }

    socket.emit("create_game", "mu nem", total_players);
    subscribe();

    
    var client = new Client();
    document.getElementById("menu").style.display = "none";
}
function connectToGame(form) {
    alert('connect to the game');
}

