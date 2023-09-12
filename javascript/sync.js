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

function main(numTasks) {
    for (let i = 0; i < numTasks; i++) {
        let item = queue[head];
        queue[head] = item + 1
        head++;
    };
};

main(numTasks);
