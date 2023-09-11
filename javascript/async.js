if (!process.argv[2] || isNaN(Number(process.argv[2], 10))) {
    throw new Error("Please provide a number of async tasks to run.");
}
let numTasks = Number(process.argv[2], 10) || 100;
if (numTasks < 100) {
    numTasks = 100
}

let queue = new Array(numTasks);
for (let i = 0; i < numTasks; i++) {
    queue[i] = 1
};
let head = 0;

async function processQueue() {
    let item = queue[head];
    queue[head] = item + 1
    head++
}

async function main(numTasks) {
    let tasks = [];
    for (let i = 0; i < numTasks; i++) {
        tasks.push(processQueue())
    };
    await Promise.all(tasks);
};

main(numTasks);
for (let i of queue) {
  console.log(i)
}
