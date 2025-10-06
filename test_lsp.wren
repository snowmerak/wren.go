// Test Wren file for LSP testing

class Calculator {
  construct new() {
    _value = 0
  }

  add(x) {
    _value = _value + x
    return _value
  }

  multiply(x) {
    _value = _value * x
    return _value
  }

  getValue() {
    return _value
  }
}

class Person {
  construct new(name, age) {
    _name = name
    _age = age
  }

  greet() {
    System.print("Hello, I'm %(_name)!")
  }

  birthday() {
    _age = _age + 1
    System.print("%(_name) is now %(_age) years old")
  }
}

var calc = Calculator.new()
calc.add(5)
calc.multiply(3)
System.print("Result: %(calc.getValue())")

var person = Person.new("Alice", 25)
person.greet()
person.birthday()
