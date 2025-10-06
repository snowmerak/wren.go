System.print("Hello from Wren CLI!")

var x = 42
System.print("The answer is %(x)")

class Point {
    construct new(x, y) {
        _x = x
        _y = y
    }
    
    distance() {
        return (_x * _x + _y * _y).sqrt
    }
}

var p = Point.new(3, 4)
System.print("Distance: %(p.distance())")
