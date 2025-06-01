import asyncio
import os
import fnmatch
from asyncio.subprocess import Process
from collections.abc import Coroutine


def get_all_modules(dirname: str) -> list[str]:
    subfolders: list[str] = [f.path for f in os.scandir(dirname) if f.is_dir()]
    for dirname in list(subfolders):
        subfolders.extend(get_all_modules(dirname))
    return subfolders


def is_module_with_tests(module: str) -> bool:
    return any([f for f in os.listdir(module) if fnmatch.fnmatch(f, "*_test.go")])


async def run_command(command: str):
    process: Process = await asyncio.create_subprocess_shell(
        command,
        stdout=asyncio.subprocess.PIPE,
        stderr=asyncio.subprocess.PIPE,
    )
    if process.stdout:
        await read_stream(process.stdout)
    if process.stderr:
        await read_stream(process.stderr)

    _ = await process.wait()


async def read_stream(stream: asyncio.StreamReader):
    while True:
        line = await stream.readline()
        if not line:
            break
        print(line.decode("utf-8").strip())


async def main():
    testComposeFile: str = "deployments/test/compose.test.yaml"
    testCommand: str = "go test -count=1 -v {module}"
    commandLine: str = "docker compose -f {composeFile}".format(
        composeFile=testComposeFile
    )

    uppers: list[Coroutine[None, None, None]] = []
    logs: list[Coroutine[None, None, None]] = []

    modules: list[str] = get_all_modules("./internal")
    for module in filter(is_module_with_tests, modules):
        moduleToTest: str = module.replace(".", "").replace("/", "-").lower()
        projectName: str = "lesta-start-test{dir}".format(dir=moduleToTest)

        mainCommand: str = commandLine + (
            " -p {projectName}".format(projectName=projectName)
        )
        upCommand: str = (
            mainCommand + " run --rm web " + testCommand.format(module=module)
        )
        logCommand: str = mainCommand + " logs -f"

        uppers.append(run_command(upCommand))
        logs.append(run_command(logCommand))

    _ = await asyncio.gather(*uppers)
    _ = await asyncio.gather(*logs)


asyncio.run(main())
