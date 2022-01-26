from __future__ import annotations
import random

from collections import deque
from dataclasses import dataclass
from pathlib import Path
from typing import Union


@dataclass
class Puzzle:
    grid: list[str]
    h: int
    w: int
    sx: int
    sy: int
    gx: int
    gy: int

    def __str__(self):
        """__str__."""
        buf = ""
        for i in range(self.h):
            row = "".join(self.grid[i])
            buf += row
            if i != self.h - 1:
                buf += "\n"
        return buf


def parse_input(file_path: Union[Path, str]) -> Puzzle:
    """parse_input from file.

    Args:
        file_path (Union[Path, str]): file_path

    Returns:
        Puzzle:
    """
    with open(file_path) as f:
        h, w = map(int, f.readline().split())
        sx, sy = map(int, f.readline().split())
        gx, gy = map(int, f.readline().split())
        g = [f.readline().rstrip() for _ in range(h)]
    return Puzzle(g, h, w, sx, sy, gx, gy)


def generate() -> Puzzle:
    """generate Puzzle

    Args:

    Returns:
        Puzzle:
    """
    h = 9
    w = 9

    sx = 0
    sy = 0
    gx = h - 1
    gy = w - 1

    sx, sy = divmod(random.randint(0, h * w - 1), w)
    gx, gy = divmod(random.randint(0, h * w - 1), w)
    while (sx, sy) == (gx, gy):
        sx, sy = divmod(random.randint(0, h * w - 1), w)
        gx, gy = divmod(random.randint(0, h * w - 1), w)
    block_num = 10
    _grid: list[list[str]] = [["."] * w for _ in range(h)]
    _grid[sx][sy] = "s"
    _grid[gx][gy] = "g"
    for _ in range(block_num):
        rx, ry = divmod(random.randint(0, h * w - 1), w)
        if (rx, ry) == (sx, sy) or (rx, ry) == (gx, gy):
            continue
        _grid[rx][ry] = "#"
    grid = ["".join(_grid[i]) for i in range(h)]

    return Puzzle(grid, h, w, sx, sy, gx, gy)


def solve(pz: Puzzle, verbose: bool = False) -> tuple[bool, int]:
    """
    Args:
        Puzzle:
        verbose:
    Returns:
        (True,3)
        (False,-1)
    if verbose, print route to the goal if possible.
    solved by BFS
    """
    g = pz.grid
    sx = pz.sx
    sy = pz.sy
    gx = pz.gx
    gy = pz.gy
    h = pz.h
    w = pz.w
    dist = [[-1] * w for _ in range(h)]
    que = deque()
    que.append((sx, sy, 0))
    dist[sx][sy] = 0

    dx = [0, 1, 0, -1]
    dy = [1, 0, -1, 0]

    def is_inside(x, y, h, w) -> bool:
        return 0 <= x and x < h and 0 <= y and y < w

    prev: dict[tuple[int, int], tuple[int, int]] = {}

    while que:
        orgx, orgy, d = que.popleft()
        x = orgx
        y = orgy
        for di in range(4):
            nx = orgx + dx[di]
            ny = orgy + dy[di]
            moved = False
            while is_inside(nx, ny, h, w) and g[nx][ny] != "#":
                x = nx
                y = ny
                nx = x + dx[di]
                ny = y + dy[di]
                moved = True
            if moved and dist[x][y] == -1:
                que.append((x, y, d + 1))
                dist[x][y] = d + 1
                prev[(x, y)] = (orgx, orgy)

    if dist[gx][gy] != -1:
        route = [(gx, gy)]
        while route[-1] != (sx, sy):
            x, y = route[-1]
            route.append(prev[route[-1]])
        route = list(reversed(route))
        if verbose:
            print(f"can reach goal by {dist[gx][gy]} step")
            print("route:", route)
        return True, dist[gx][gy]
    else:
        if verbose:
            print("cannot reach")
        return False, -1


def save_puzzle(save_name: str, pz: Puzzle):
    with open("resources/" + save_name, "w") as f:
        f.write(f"{pz.h} {pz.w}\n")
        f.write(f"{pz.sx} {pz.sy}\n")
        f.write(f"{pz.gx} {pz.gy}\n")
        for row in pz.grid:
            f.write(row + "\n")


def main():
    # generate puzzles solvable
    MAX_PZ = 10
    for i in range(MAX_PZ):
        pz = generate()
        ok, step = solve(pz)
        while not ok:
            pz = generate()
            ok, step = solve(pz)
        save_puzzle(f"{i}.txt", pz)


if __name__ == "__main__":
    main()
