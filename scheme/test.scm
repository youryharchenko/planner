(define (square x) (* x x))

(define (average x y)
  (/ (+ x y) 2)
)

(define (sum term a next b)
  (if (> a b)
    0
    (+ (term a)
      (sum term (next a) next b)
    )
  )
)

(define (integral f a b dx)
  (define (add-dx x)
    (+ x dx)
  )
  (*
    (sum
      f
      (+ a (/ dx 2.0))
      add-dx
      b
    )
    dx
  )
)

(define (inc n) (+ n 1))

(define (sum-cubes a b)
  (sum cube a inc b)
)

(define (cube x) (* x x x))

(define (identity x) x)

(define (sum-integers a b)
  (sum identity a inc b)
)

(define (pi-sum a b)
  (define (pi-term x)
    (/ 1.0 (* x (+ x 2))))
  (define (pi-next x)
    (+ x 4))
    (sum pi-term a pi-next b))

(define tolerance 0.00001)

(define (fixed-point f first-guess)
  (define (close-enough? v1 v2)
    (< (abs (- v1 v2)) tolerance))
  (define (try guess)
    (let ((next (f guess)))
      (if (close-enough? guess next)
        next
        (try next))))
  (try first-guess))

(define (average-damp f)
  (lambda (x) (average x (f x))))

(define (cube-root x)
  (fixed-point (average-damp (lambda (y) (/ x (square y)))) 1.0))

(define (deriv g)
  (lambda (x)
  (/ (- (g (+ x dx)) (g x))
  dx)))

(define (newton-transform g)
  (lambda (x)
    (- x (/ (g x) ((deriv g) x)))))

(define (newtons-method g guess)
    (fixed-point (newton-transform g) guess))

(define (sqrt x)
    (newtons-method (lambda (y) (- (square y) x)) 1.0))

(define (add-rat x y)
    (make-rat (+ (* (numer x) (denom y))
    (* (numer y) (denom x)))
    (* (denom x) (denom y))))
(define (sub-rat x y)
    (make-rat (- (* (numer x) (denom y))
    (* (numer y) (denom x)))
    (* (denom x) (denom y))))
(define (mul-rat x y)
    (make-rat (* (numer x) (numer y))
    (* (denom x) (denom y))))
(define (div-rat x y)
    (make-rat (* (numer x) (denom y))
    (* (denom x) (numer y))))
(define (equal-rat? x y)
    (= (* (numer x) (denom y))
    (* (numer y) (denom x))))

(define (gcd a b)
  (if (= b 0)
    a
    (gcd b (remainder a b))))

(define (cons x y)
  (define (dispatch m)
    (cond ((= m 0) x)
    ((= m 1) y)
    (else (error "Argument not 0 or 1 - CONS" m))))
  dispatch)
(define (car z) (z 0))
(define (cdr z) (z 1))

(define (list-ref items n)
  (if (= n 0)
    (car items)
    (list-ref (cdr items) (- n 1))))
