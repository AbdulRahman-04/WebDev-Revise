// console.log("Hello World in js🤌!");
// console.log("hello world 2");
// console.log("hello world 3");

// Variables in js

/*
  Variables in javascript are like containers which stores data in itself
  there are three ways t odeclare a variable in js :

    1. var
    2. let
    3. const
    4. default

    var : es5 method, can be re-declared and updated 
    let : es6 method, can be updated but cannot be re-declared
    const : es6 method, cannot be updated or re-declared (constant value)

*/

var myVar = "Rahman";
// console.log(myVar);

var myVar = "Hello";
console.log(myVar);

let name = "syed bhai";
name = "rahman bhai";

console.log(name);

const myName = "Syed Abdul Rahman";
console.log(myName);

rollNo = 5049;
console.log(rollNo);

// here roll no variable is using var keyword behind the scenes.

// Data types in js

/*

Primtive data types :
1. String
2 Number
3. Boolean 
4. null
5. undefined
6. symbol
7. BigInt (big integer numbers)

Non primitive data types :
1. Array 2. Object 3.Function

*/

// Scoping in variables:

// global scope:                VAR                LET             CONST

// inside of block:              ✅                ✅               ✅
// outside of block:              ✅                ✅               ✅

// FUNCTION SCOPE:              VAR                LET             CONST

// inside of block:              ✅                 ✅              ✅
// outside of block:              ❌                ❌               ❌

// block scope:                  VAR                LET             CONST

// inside of block:               ✅                 ✅              ✅
// outside of block:               ✅                ❌               ❌

// ==========================
// GLOBAL SCOPE
// ==========================
var x = 10; // function scoped or global
let y = 20; // block scoped
const z = 12; // block scoped, constant

console.log("Global var x:", x); // 10
console.log("Global let y:", y); // 20
console.log("Global const z:", z); // 12

// ==========================
// FUNCTION SCOPE
// ==========================
function myTest() {
  var x = 12; // function scoped
  console.log("Function myTest var x:", x); // 12
}
myTest();

var x = 19;
console.log("Global var x after myTest:", x); // 19

function myTest2() {
  let y = 29; // function/block scoped
  console.log("Function myTest2 let y:", y); // 29
}
myTest2();

// console.log(y); // Error: y is not defined (let is block scoped)

function myTest3() {
  const z = 49; // function/block scoped
  console.log("Function myTest3 const z:", z); // 49
}
myTest3();

console.log("Global const z:", z); // 12

// ==========================
// BLOCK SCOPE
// ==========================
if (true) {
  var a = 100; // var ignores block scope → behaves like global
  let b = 200; // block scoped
  const c = 300; // block scoped

  console.log("Inside block var a:", a); // 100
  console.log("Inside block let b:", b); // 200
  console.log("Inside block const c:", c); // 300
}

console.log("Outside block var a:", a); // 100
// console.log(b); // Error: b is not defined
// console.log(c); // Error: c is not defined

// ==========================
// LOOP BLOCK SCOPE
// ==========================
for (let i = 0; i < 3; i++) {
  console.log("Inside loop i:", i); // 0,1,2
}
// console.log(i); // Error: i is not defined

/*
 matlab bro basically var keyword variables ku har jgah
 access nd redeclre/reassign krskte but let nd const variables ku jaha declare
 kre wahi access krskte like function or blocks
*/

// Type conversion in js : converting a variable value of one data type to another.

// implicit conversion || explicit conversion : 

// implicit conversion : done by js engine automatically. 

// (+) converts number || boolean to string (also called string concatination)
let ab = "1" + 9;
console.log(ab, typeof ab);



// (-) converts string to number
let ac = "5" - 3;
console.log(ac, typeof ac);


// Explicit conversion: done by developer manually using methods.

// Conversion to String: 

  let var1 = true
  console.log(String(var1), typeof String(var1));


// Conversion to Number:

let var2 = false
console.log(Number(var2), typeof Number(var2));


// conversion to boolean 

let var3 = 0
console.log(Boolean(var3), typeof Boolean(var3));

