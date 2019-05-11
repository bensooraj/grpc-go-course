let PROTO_PATH = "../../calculator/calculatorpb";
let PROTO_FILE = "../../calculator/calculatorpb/calculator.proto";

const protoLoader = require('@grpc/proto-loader');
const grpc = require('grpc');

const packageDefinition = protoLoader.loadSync(PROTO_FILE, {
    keepCase: true,
    includeDirs: [PROTO_PATH],
    longs: Number
});

const calProto = grpc.loadPackageDefinition(packageDefinition).calculator;

// ########################################################################

async function main() {
    const client = new calProto.CalculatorService(
        "localhost:50051",
        grpc.credentials.createInsecure()
    );

    // Unary
    // try {
    //     console.log("Unary: Sum")
    //     await sum(client);
    //     console.log("########################################");
    //     console.log();
    // } catch (error) {
    //     console.log("Sum error: ", error);
    // }

    // Server Streaming
    // try {
    //     console.log("Server Streaming: PrimeNumberDecomposition")
    //     await primeNumberDecomposition(client);
    //     console.log("########################################");
    //     console.log();
    // } catch (error) {
    //     console.log("primeNumberDecomposition error: ", error);
    // }

    // Client Streaming
    try {
        console.log("Client Streaming: ComputeAverage")
        await computeAverage(client);
        console.log("########################################");
        console.log();
    } catch (error) {
        console.log("computeAverage error: ", error);
    }

    // BiDi Streaming
    // try {
    //     console.log("BiDi Streaming: FindMaximum")
    //     await findMaximum(client);
    //     console.log("########################################");
    //     console.log();
    // } catch (error) {
    //     console.log("findMaximum error: ", error);
    // }
}
// 
main();

function sum(client) {
    return new Promise((resolve, reject) => {
        client.Sum({
            first_number: 123,
            second_number: 456
        }, function (err, response) {
            if (err) {
                reject(err);
                return;
            }
            console.log("response", response);
            resolve(response);
            return;
        })
    });
}

function primeNumberDecomposition(client) {
    return new Promise((resolve, reject) => {
        // 
        const stream = client.PrimeNumberDecomposition({
            number: 444
        });
        // 
        stream.on('data', function (factor) {
            console.log("factor: ", factor);
        });
        // 
        stream.on('end', resolve);
    });
}

function findMaximum(client) {
    return new Promise(async (resolve, reject) => {
        // Get the stream object
        const stream = client.FindMaximum({
            number: 444
        });
        stream.on('data', function (maximum) {
            console.warn("[RECEIVING] Current maximum: ", maximum);
        });
        // 
        stream.on('end', resolve);

        // 
        const numberArray = [1, 4, 5, 3, 35, 5, 1, 25, 6, 45, 23, 3, 2, 5, 56];
        for (const number of numberArray) {
            await sleep(1000);
            stream.write({
                number
            });
            console.log("[SENDING] Number: ", number);
        }
        stream.end()
    });
}

function computeAverage(client) {
    return new Promise(async (resolve, reject) => {
        // Get the stream object
        let stream = client.ComputeAverage({
            number: 1,
        }, (err, response) => {
            if (err) {
                console.log("computeAverage | client.ComputeAverage | Error: ", err);
                reject(err);
                return;
            }
            console.log("computeAverage | client.ComputeAverage | Response: ", response);
            resolve();
            return;
        });

        const numberArray = [1, 4, 5, 3, 35, 5, 1, 25, 6, 45, 23, 3, 2, 5, 56];
        for (const number of numberArray) {
            await sleep(100);
            stream.write({
                number
            });
            console.log("[SENDING] Number: ", number);
        }
        stream.end()
    });
}

async function sleep(ms) {
    return new Promise(resolve => {
        setTimeout(resolve, ms)
    })
}