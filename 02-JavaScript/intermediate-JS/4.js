// Arrays : an array is also like object, which can store multiple datatype value initself.

// const myArr = ["hi", 12, false, null, undefined]

// console.log(myArr);

// Crud operations on array:

// updating
// myArr[0] = "Hello bro"
// console.log(myArr);

// myArr[4] = true;
// console.log(myArr);


// deleteing
// delete myArr[2] 
// delete myArr[3]
// console.log(myArr);


// Array methods: 

// 1.at(): it shows the which element is present at the given index value

let myArr = [1, true, 'hey', 'lol']
console.log(myArr.at(1));

// 2.Concat(): it merges two array elements 
let arr1 = [1, 2];
let arr2 = [true, false]
console.log(arr1.concat(arr2));

// 3.fill: it fills a value provided by us in the given index position.

let myArr1 = [ , 'bruv',  , 'kaha ho?']
console.log(myArr1.fill('hi', 0., 1));
console.log(myArr1.fill('tum', 2, 3));

// 4.filter(): this array method filters out the elements which satisfies the condition inside arrow function

let myArr3 = [20, 40, 60, 80, 100]

console.log(myArr3.filter((x)=>{
    return x > 50
}));


// 5.find: it finds the value of index position provided and returns its complete value and not its index

let myArr4 = ['banana', 'yooo', 'jojo', 'cr9', 'goDaddy'];

console.log(myArr4.find((x)=>{
    return x === "jojo"
}));


// 6.flat(): it converts back the sub arrays stored inside main array into main array.

let myArr5 = [[35], [89], 88];
console.log(myArr5.flat());


// 7.forEach(): it checks all the elements inside of an array and prints their index value but it can't 
//               return the value and also can't use more array methods inside it.

// “map/forEach: first param = value, second = index.” ✅

let myArr6 = ['hi', 'bolo', false, 99, 'bolo'];

console.log(myArr6.forEach((index, value)=>{
    console.log(value, index);
}));


// 8.includes: it checks if your given value is included in the array or not

let myArr8 = ['samid', 'suhail', 'saad', 'syed', 'mz'];
console.log(myArr8.includes('saad'));


// 9.indexof: it also returns the index of given element from the array but returns -1 if not found.

let myArr9 = [29, 354, true, false]
console.log(myArr9.indexOf('samid'));


// 10.isArray(): it checks if the given value is array or not.

let myArr10 = [19, 22, 56, 88];
console.log(Array.isArray(myArr10));

// 11.join(): it joins the array elements with a given symbol inside join method
let myArr11 = ['ab', 'smd', 'ism', 'nasir'];
console.log(myArr11.join('-'));

// 12.map(): it is also exact same like for each but map can also return the values of array and can
//           perform different array methods simultaniously

// value → current element ka value (ye sabse pehle aata)

// index → current element ka index (0,1,2…)

// “map/forEach: first param = value, second = index.” ✅


let myArr12 = ['haan', ' bolo', ' kya ', 'hua', 'bhaiyya'];
console.log(myArr12.map((value, index)=>{
    return `${index}: ${value}`
}));


const numbers = [1, 2, 3, 4, 5];

console.log(numbers.map((value, index)=>{
    return `${index}: ${value*2}`
}));

const fruits = ['apple', 'banana', 'mango'];
console.log(fruits.map((value, index)=> {
   
    return `${index}: ${value.toLocaleUpperCase()}`

}));

// 13.pop(): this array method removes the last element of an array
let myArr14 = ['han', 'kidr', 'h', 'bey?', 'b*astard'];
myArr14.pop()
console.log(myArr14);


// 14.push(): this array method pushes a new element to the last index of array
let myArr15 = ['han', 'kidr', 'h', 'bey?', 'b*astard'];
myArr15.push("?")
console.log(myArr15);

// 15.reduce(): it is used to add sum of all array elements
let myArr16 = [20, 40, 60, 80, 100];
console.log(myArr16.reduce((acc, curr)=> {
 return acc+curr
}, 0));



// 16.reverse(): 
let myArr17 = ['g', 'k', 'l ', 'o', 'q', 'z'];
console.log(myArr17.reverse());

// 17.shift : it removes the first array index element
let myArr18= [20, 40, 60, 80, 100];
myArr18.shift();
console.log(myArr18);


// 18.slice: it cuts a piece from an array from the given index value to an actual length of array
let myArr19 = [20, 40, 60, 80, 100];
console.log(myArr19.slice(2, 3));

// 19.some: it serches the array element just like find but returns the output in true or false and doesn't
//          stores the actual value.

let myArr20 = [20, 40, 60, 80, 100];
console.log(myArr20.some((x)=> {
    return x === 80
}));


// 20.splice: it removes the specific element and new element in place of it.
let myArr21 = [20, 40, 60, 80, 100, 120, 140, 160];
// console.log(myArr21.splice(1, 3, 50, 70));


// 21.toString(): converts the whole array into string.
let myArr22 = [true, false, 45, 'hey']
console.log(myArr22.toString());


// 22.unshift(): it adds a new element at first index of array:
let myArr23= ['bol', true, 'hai', 'ya', false];
myArr23.unshift("bhai, ");
console.log(myArr23);



// 23.length(): it checks actual length of array and not its  index value
// nai aara ye ek baar suhail bhai se pucho


// 24.sort(): it sorts the array element in an order.
let thisArr = ['g', 'i', 'a', 'k', 'm', 'n', 'j'];
console.log(thisArr.sort());

// 25.tolocalestring(): this also converts array elements into string

