from __future__ import annotations

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


def parse_input(file_path: Union[Path, str]) -> Puzzle:
    with open(file_path) as f:
        h, w = map(int, f.readline().split())
        sx, sy = map(int, f.readline().split())
        gx, gy = map(int, f.readline().split())
        g = [f.readline().rstrip() for _ in range(h)]
    return Puzzle(g, h, w, sx, sy, gx, gy)


def solve(pz: Puzzle) -> bool:
    """
    solved by BFS
    print route to the goal if possible
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
        print(f"can reach goal by {dist[-1][-1]} step")
        print("route:", route)
        return True
    else:
        print("cannot reach")
        return False


def main():
    pz: Puzzle = parse_input(Path("./resources/sample_layer.txt"))
    print(*pz.grid,sep='\n')
    solve(pz)


if __name__ == "__main__":
    main()
