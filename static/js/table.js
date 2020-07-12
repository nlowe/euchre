const tableName = document.head.querySelector('meta[name="x-table-name"]').content;
const clientID = uuidv4()
const seen = new Set();

console.info("I am", clientID, "in table", tableName)

const qr = qrcode(0, 'M');
qr.addData(window.location.href);
qr.make();
document.getElementById('join-code').innerHTML = qr.createSvgTag();


const proto = window.location.protocol === 'http:' ? 'ws' : 'wss';
const path = window.location.origin + '/_play';
const address = proto + path.substr(path.indexOf(':'), path.length);
console.info("Connecting to table", tableName, " at ", address)

const connection = new autobahn.Connection({
    url: address,
    realm: tableName
});

connection.onopen = function(session) {
    session.subscribe('state', function (args, kwargs) {
        console.info('Got new state', args, kwargs)
        const from = args[0]
        if (from === clientID) {
            console.info("Not saying hello to ourself")
            return
        }

        if (seen.has(from)) {
            console.info("Already said hi to", from)
            return
        }

        console.info("Saying hello to", from)
        seen.add(from)
        session.call('say', [clientID, 'hello', from]).then(function (res) {
            console.info("RPC 'say' complete:", res)
        })
    })

    session.call('say', [clientID, 'hello!']).then(function (res) {
        console.info("RPC 'say' complete:", res)
    })
}

connection.open()
