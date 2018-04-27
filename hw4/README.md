# HW 4

## Description
For this homework I did not solve the 2-point boundary value problem. I instead opted for forward simulation. This was done using with the simplest numerical method I know, namely Eulers method. I then forward simulated in the direction of the sampled point over a distance of 𝛅. The produced path were then check for point-wise feasibility of the constraints for position and velocity constraints. The accelerations were computed and then clamped to also adhere to its constraints. Note that the RRT shown in the figures are only showing straight lines between the positions of the endpoints, although there are some feasible edge there, but that cluttered up stuff waaay too much to show.

# Usage
```shell
# for problem 1 in problems.json
go run config.go main.go rrt.go -p 1 | python plot.py

# more usage help
go run config.go main.go rrt.go -h
```