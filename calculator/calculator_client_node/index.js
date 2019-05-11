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
    try {
        console.log("Unary: Sum")
        await sum(client);
        console.log("########################################");
        console.log();
    } catch (error) {
        console.log("Sum error: ", error);
    }

    // Server Streaming
    try {
        console.log("Server Streaming: PrimeNumberDecomposition")
        await primeNumberDecomposition(client);
        console.log("########################################");
        console.log();
    } catch (error) {
        console.log("primeNumberDecomposition error: ", error);
    }
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