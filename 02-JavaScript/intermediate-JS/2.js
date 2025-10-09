// Operators in js :

// there are several types of operators in js :

// 1. Arithmetic operators : +, -, *, /, %, ++, --

var a = 10;
var b = 20;

console.log(a + b, a - b, a * b, a % b, ++a, --b);

// 2. Assignment operators : =, +=, -=, *=, /=, %=

var c = 35;

c += 5;
console.log(c);

var d = 50;

d %= 30;
console.log(d);

// Comparison operators : ==, ===, !=, !==, >, <, >=, <=

let ab = "123";
let ba = 123;

//    console.log(ab == ba);
//    console.log(ab === ba);
//    console.log(ab != ba);
//    console.log(ab !== ba);

// Logical operators : &&, ||, !
let x = 10;
let y = 20;

//   console.log(x <= y && y >= x);
console.log(x >= y || y > x);
console.log(!(x >y) );


// Ternary Operator: used for single line conditions.

let age = 17;

age >= 18? console.log("You can drive") : console.log("You are a minor, u cant drive");

