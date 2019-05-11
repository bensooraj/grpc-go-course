let PROTO_PATH = "../../calculator/calculatorpb";
let PROTO_FILE = "../../calculator/calculatorpb/calculator.proto";

const protoLoader = require('@grpc/proto-loader');
const grpc = require('grpc');

const packageDefinition = protoLoader.loadSync(PROTO_FILE, {
    keepCase: true,
    includeDirs: [PROTO_PATH]
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
        await sum(client);
    } catch (error) {
        console.log("Sum error: ", error);
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