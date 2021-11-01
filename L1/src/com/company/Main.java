package com.company;

public class Main {

    public static void main(String[] args) {
        Board board = Board.getInstance();
        Board.printBoard(board.getBoardContent());

        long start = System.currentTimeMillis();
        Solver solver = new Solver();

        long end = System.currentTimeMillis();
        System.out.println("time taken " + (end-start) + " milli seconds");


    }
}
