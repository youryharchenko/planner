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
