if (!process.argv[2] || isNaN(Number(process.argv[2], 10))) {
    throw new Error("Please provide a number of async tasks to run.");
}
let numTasks = Number(process.argv[2], 10) || 100;
if (numTasks < 100) {
    numTasks = 100
}

let queue = new Array(1000 + numTasks);
for (let i = 0; i < 1000; i++) {
    queue[i] = 1
};
let head = 0;
let tail = 1000;

function main(numTasks) {
    for (let i = 0; i < numTasks; i++) {
        let item = queue[head];
        head++;
        queue[tail] = item + 1;
        tail++
    };
};

main(numTasks);
console.log(head);
console.log(queue.length)
