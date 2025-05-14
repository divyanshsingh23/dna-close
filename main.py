#!/usr/bin/env python3

import os
import sys
import math
from typing import List, Tuple
from dataclasses import dataclass
import heapq

@dataclass
class Sample:
    name: str
    coordinates: List[float]
    distance: float = 0.0

def clear_screen():
    """Clear the terminal screen."""
    os.system('clear' if os.name != 'nt' else 'cls')

def read_xyz_coordinates(filename: str) -> List[float]:
    """Read G25 coordinates from xyz.txt file."""
    try:
        with open(filename, 'r') as f:
            line = f.readline().strip()
            coords = [float(x) for x in line.split(',')]
            if len(coords) != 25:
                raise ValueError(f"Expected 25 coordinates, found {len(coords)}")
            return coords
    except FileNotFoundError:
        print(f"\nError: {filename} not found.")
        print("Please create this file with your 25 G25 coordinates (comma-separated).")
        sys.exit(1)
    except (ValueError, IndexError) as e:
        print(f"\nError: Invalid format in {filename}")
        print("File should contain exactly 25 comma-separated numbers.")
        sys.exit(1)

def euclidean_distance(coords1: List[float], coords2: List[float]) -> float:
    """Calculate Euclidean distance between two coordinate sets."""
    return math.sqrt(sum((a - b) ** 2 for a, b in zip(coords1, coords2)))

def find_closest_matches(user_coords: List[float], samples_file: str, top_n: int = 10) -> List[Sample]:
    """Find the closest matching samples using Euclidean distance."""
    closest = []  # Will store (distance, name) tuples
    
    try:
        with open(samples_file, 'r') as f:
            for line in f:
                line = line.strip()
                if not line:
                    continue
                    
                try:
                    parts = line.split(',')
                    name = parts[0]
                    coords = [float(x) for x in parts[1:26]]  # 25 coordinates
                    
                    if len(coords) != 25:
                        print(f"Warning: Skipping sample {name} - invalid coordinate count")
                        continue
                        
                    distance = euclidean_distance(user_coords, coords)
                    sample = Sample(name, coords, distance)
                    
                    # Use a max heap to keep track of top_n closest matches
                    if len(closest) < top_n:
                        heapq.heappush(closest, (-distance, sample))
                    else:
                        heapq.heappushpop(closest, (-distance, sample))
                        
                except (ValueError, IndexError):
                    print(f"Warning: Skipping malformed line for sample {parts[0] if parts else 'unknown'}")
                    continue
                    
    except FileNotFoundError:
        print(f"\nError: {samples_file} not found.")
        print("Please ensure the samples file exists in the current directory.")
        sys.exit(1)
    
    # Convert to list of samples, sorted by distance
    return [sample for _, sample in sorted(closest, key=lambda x: -x[0])]

def count_samples(filename: str) -> int:
    """Count the number of valid samples in the samples file."""
    try:
        with open(filename, 'r') as f:
            return sum(1 for line in f if line.strip())
    except FileNotFoundError:
        print(f"\nError: {filename} not found.")
        sys.exit(1)

def display_menu():
    """Display the main menu."""
    print("\nG25 Coordinate Matcher")
    print("=====================")
    print("[1] Compare my coordinates")
    print("[2] Show sample count")
    print("[3] Exit")
    return input("\nSelect an option (1-3): ").strip()

def main():
    """Main program loop."""
    while True:
        clear_screen()
        choice = display_menu()
        
        if choice == '1':
            clear_screen()
            print("Comparing coordinates...")
            try:
                user_coords = read_xyz_coordinates('xyz.txt')
                matches = find_closest_matches(user_coords, 'samples.txt')
                
                print("\nTop 10 Closest Matches:")
                print("=======================")
                for i, match in enumerate(matches, 1):
                    print(f"{i}. {match.name:<30} Distance: {match.distance:.6f}")
                
            except Exception as e:
                print(f"\nAn error occurred: {str(e)}")
            
            input("\nPress Enter to continue...")
            
        elif choice == '2':
            clear_screen()
            count = count_samples('samples.txt')
            print(f"\nTotal number of samples: {count}")
            input("\nPress Enter to continue...")
            
        elif choice == '3':
            clear_screen()
            print("Thank you for using G25 Coordinate Matcher!")
            sys.exit(0)
            
        else:
            print("\nInvalid option. Please try again.")
            input("Press Enter to continue...")

if __name__ == "__main__":
    main()
