struct waiter {
    int f, n;
    std::mutex m;
    dstd::condition_variable c;
    waiter(int n) : f(0), n(n) {}
    void step() {
        n--;
    }
    void signal() {
        m.lock();
        f = 1;
        m.unlock();
        c.notify_one();
    }
    int wait() {
        if (n) {
            std::unique_lock<std::mutex> lk(m);
            if (!f) c.wait(lk);
            f = 0;
            return 1;
        }
        return 0;
    }
};

class ZeroEvenOdd {
private:
    waiter wz, w1, w2;
public:
    ZeroEvenOdd(int n) : wz(n), w1((n+1)/2), w2(n/2) { wz.signal(); }

    // printNumber(x) outputs "x", where x is an integer.
    void zero(function<void(int)> printNumber) {   
        waiter& w = wz;
        for (int s = 0; w.wait(); w.step(), s++) {
            printNumber(0);
            (s % 2 == 0 ? w1 : w2).signal();
        }
    }
    void even(function<void(int)> printNumber) {   
        waiter& w = w2;
        for (int s = 2; w.wait(); w.step(), s += 2) {
            printNumber(s);
            wz.signal();
        }
    }
    void odd(function<void(int)> printNumber) {   
        waiter& w = w1;
        for (int s = 1; w.wait(); w.step(), s += 2) {
            printNumber(s);
            wz.signal();
        }
    }
};