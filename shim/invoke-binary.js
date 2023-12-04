const childProcess = require('child_process')
const os = require('os')
const process = require('process')

function chooseBinary() {
    const platform = os.platform()
    const arch = os.arch()

    if (platform === 'linux' && arch === 'x64') {
        return `assemble-linux-amd64`
    }
    if (platform === 'linux' && arch === 'arm64') {
        return `assemble-linux-arm64`
    }
    if (platform === 'windows' && arch === 'x64') {
        return `assemble-windows-amd64`
    }
    if (platform === 'windows' && arch === 'arm64') {
        return `assemble-windows-arm64`
    }
    if (platform === 'darwin' && arch === 'x64') {
        return `assemble-darwin-amd64`
    }
    if (platform === 'darwin' && arch === 'arm64') {
        return `assemble-darwin-arm64`
    }


    console.error(`Unsupported platform (${platform}) and architecture (${arch})`)
    process.exit(1)
}

function main() {
    const binary = chooseBinary()
    const mainScript = `${__dirname}/../bin/${binary}`

    console.log("Spawning %s", mainScript)
    const spawnSyncReturns = childProcess.spawnSync(mainScript, { stdio: 'inherit' })
    const status = spawnSyncReturns.status
    console.log("Captured status %s", status)
    if (typeof status === 'number') {
        process.exit(status)
    }
    process.exit(1)
}

if (require.main === module) {
    main()
}
